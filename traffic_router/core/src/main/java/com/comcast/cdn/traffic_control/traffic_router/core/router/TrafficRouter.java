/*
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import java.io.IOException;
import java.net.InetAddress;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;
import java.util.stream.Collectors;

import com.comcast.cdn.traffic_control.traffic_router.configuration.ConfigurationListener;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringResult;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringTarget;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.Steering;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringRegistry;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.ConsistentHasher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.MaxmindGeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.log4j.Logger;
import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.xbill.DNS.Name;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation.LocalizationMethod;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.ZoneManager;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.DNSAccessRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringGeolocationComparator;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationException;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNodeException;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeo;
import com.comcast.cdn.traffic_control.traffic_router.core.request.DNSRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.RouteType;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIp;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIpDatabaseService;

public class TrafficRouter {
	public static final Logger LOGGER = Logger.getLogger(TrafficRouter.class);
	public static final String XTC_STEERING_OPTION = "x-tc-steering-option";

	private final CacheRegister cacheRegister;
	private final ZoneManager zoneManager;
	private final GeolocationService geolocationService;
	private final GeolocationService geolocationService6;
	private final AnonymousIpDatabaseService anonymousIpService;
	private final FederationRegistry federationRegistry;
	private final boolean consistentDNSRouting;

	private final Random random = new Random(System.nanoTime());
	private Set<String> requestHeaders = new HashSet<String>();
	private static final Geolocation GEO_ZERO_ZERO = new Geolocation(0,0);
	private ApplicationContext applicationContext;

	private final ConsistentHasher consistentHasher = new ConsistentHasher();
	private SteeringRegistry steeringRegistry;

	private final Map<String, Geolocation> defaultGeolocationsOverride = new HashMap<String, Geolocation>();

	public TrafficRouter(final CacheRegister cr, 
			final GeolocationService geolocationService, 
			final GeolocationService geolocationService6,
			final AnonymousIpDatabaseService anonymousIpService,
			final StatTracker statTracker,
			final TrafficOpsUtils trafficOpsUtils,
			final FederationRegistry federationRegistry,
			final TrafficRouterManager trafficRouterManager) throws IOException {
		this.cacheRegister = cr;
		this.geolocationService = geolocationService;
		this.geolocationService6 = geolocationService6;
		this.anonymousIpService = anonymousIpService;
		this.federationRegistry = federationRegistry;
		this.consistentDNSRouting = JsonUtils.optBoolean(cr.getConfig(), "consistent.dns.routing");
		this.zoneManager = new ZoneManager(this, statTracker, trafficOpsUtils, trafficRouterManager);

		if (cr.getConfig() != null) {
			// maxmindDefaultOverride: {countryCode: , lat: , long: }
			final JsonNode geolocations = cr.getConfig().get("maxmindDefaultOverride");
			if (geolocations != null) {
				for (final JsonNode geolocation : geolocations) {
					final String countryCode = JsonUtils.optString(geolocation, "countryCode");
					final double lat = JsonUtils.optDouble(geolocation, "lat");
					final double longitude = JsonUtils.optDouble(geolocation, "long");
					defaultGeolocationsOverride.put(countryCode, new Geolocation(lat, longitude));
				}
			}
		}
	}

	public ZoneManager getZoneManager() {
		return zoneManager;
	}

	/**
	 * Returns a {@link List} of all of the online {@link Cache}s that support the specified
	 * {@link DeliveryService}. If no online caches are found to support the specified
	 * DeliveryService an empty list is returned.
	 * 
	 * @param ds
	 *            the DeliveryService to check
	 * @return collection of supported caches
	 */
	public List<Cache> getSupportingCaches(final List<Cache> caches, final DeliveryService ds) {
		final List<Cache> supportingCaches = new ArrayList<Cache>();

		for (final Cache cache : caches) {
			if (!cache.hasDeliveryService(ds.getId())) {
				continue;
			}

			if (cache.hasAuthority() ? cache.isAvailable() : true) {
				supportingCaches.add(cache);
			}
		}

		return supportingCaches;
	}

	public CacheRegister getCacheRegister() {
		return cacheRegister;
	}
	protected DeliveryService selectDeliveryService(final Request request, final boolean isHttp) {
		if(cacheRegister==null) {
			LOGGER.warn("no caches yet");
			return null;
		}

		final DeliveryService deliveryService = cacheRegister.getDeliveryService(request, isHttp);

		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug("Selected DeliveryService: " + deliveryService);
		}
		return deliveryService;
	}

	boolean setState(final JsonNode states) throws UnknownHostException {
		setCacheStates(states.get("caches"));
		setDsStates(states.get("deliveryServices"));
		return true;
	}
	private boolean setDsStates(final JsonNode dsStates) {
		if(dsStates == null) {
			return false;
		}
		final Map<String, DeliveryService> dsMap = cacheRegister.getDeliveryServices();
		for (final String dsName : dsMap.keySet()) {
			dsMap.get(dsName).setState(dsStates.get(dsName));
		}
		return true;
	}
	private boolean setCacheStates(final JsonNode cacheStates) {
		if(cacheStates == null) {
			return false;
		}
		final Map<String, Cache> cacheMap = cacheRegister.getCacheMap();
		if(cacheMap == null) { return false; }
		for (final String cacheName : cacheMap.keySet()) {
			final String monitorCacheName = cacheName.replaceFirst("@.*", "");
			final JsonNode state = cacheStates.get(monitorCacheName);
			cacheMap.get(cacheName).setState(state);
		}
		return true;
	}

	protected static final String UNABLE_TO_ROUTE_REQUEST = "Unable to route request.";
	protected static final String URL_ERR_STR = "Unable to create URL.";

	public GeolocationService getGeolocationService() {
		return geolocationService;
	}

	public AnonymousIpDatabaseService getAnonymousIpDatabaseService() {
		return anonymousIpService;
	}

	public Geolocation getLocation(final String clientIP) throws GeolocationException {
		return clientIP.contains(":") ? geolocationService6.location(clientIP) : geolocationService.location(clientIP);
	}

	private GeolocationService getGeolocationService(final String geolocationProvider, final String deliveryServiceId) {
		if (applicationContext == null) {
			LOGGER.error("ApplicationContext not set unable to use custom geolocation service providers");
			return null;
		}

		if (geolocationProvider == null || geolocationProvider.isEmpty()) {
			return null;
		}

		try {
			return (GeolocationService) applicationContext.getBean(geolocationProvider);
		} catch (Exception e) {
			StringBuilder error = new StringBuilder("Failed getting providing class '" + geolocationProvider + "' for geolocation");
			if (deliveryServiceId != null && !deliveryServiceId.isEmpty()) {
				error = error.append(" for delivery service " + deliveryServiceId);
			}
			error = error.append(" falling back to " + MaxmindGeolocationService.class.getSimpleName());
			LOGGER.error(error);
		}

		return null;
	}

	public Geolocation getLocation(final String clientIP, final String geolocationProvider, final String deliveryServiceId) throws GeolocationException {
		final GeolocationService customGeolocationService = getGeolocationService(geolocationProvider, deliveryServiceId);
		return customGeolocationService != null ? customGeolocationService.location(clientIP) : getLocation(clientIP);
	}

	public Geolocation getLocation(final String clientIP, final DeliveryService deliveryService) throws GeolocationException {
		return getLocation(clientIP, deliveryService.getGeolocationProvider(), deliveryService.getId());
	}

	public List<Cache> getCachesByGeo(final DeliveryService ds, final Geolocation clientLocation, final Track track) throws GeolocationException {
		int locationsTested = 0;

		final int locationLimit = ds.getLocationLimit();
		final List<CacheLocation> geoEnabledCacheLocations = filterEnabledLocations(getCacheRegister().getCacheLocations(), LocalizationMethod.GEO);
		final List<CacheLocation> cacheLocations1 = ds.filterAvailableLocations(geoEnabledCacheLocations);
		final List<CacheLocation> cacheLocations = orderCacheLocations(cacheLocations1, clientLocation);

		for (final CacheLocation location : cacheLocations) {
			final List<Cache> caches = selectCaches(location, ds);
			if (caches != null) {
				track.setResultLocation(location.getGeolocation());
				if (track.getResultLocation().equals(GEO_ZERO_ZERO)) {
					LOGGER.error("Location " + location.getId() + " has Geolocation " + location.getGeolocation());
				}
				return caches;
			}
			locationsTested++;
			if(locationLimit != 0 && locationsTested >= locationLimit) {
				return null;
			}
		}

		return null;
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	protected List<Cache> selectCaches(final HTTPRequest request, final DeliveryService ds, final Track track) throws GeolocationException {
		CacheLocation cacheLocation;
		ResultType result = ResultType.CZ;
		final boolean useDeep = (ds.getDeepCache() == DeliveryService.DeepCachingType.ALWAYS);

		if (useDeep) {
			// Deep caching is enabled. See if there are deep caches available
			cacheLocation = getDeepCoverageZoneCacheLocation(request.getClientIP(), ds);
			if (cacheLocation != null && cacheLocation.getCaches().size() != 0) {
				// Found deep caches for this client, and there are caches that might be available there.
				result = ResultType.DEEP_CZ;
			} else {
				// No deep caches for this client, would have used them if there were any. Fallback to regular CZ
				cacheLocation = getCoverageZoneCacheLocation(request.getClientIP(), ds);
			}
		} else {
			// Deep caching not enabled for this Delivery Service; use the regular CZ
			cacheLocation = getCoverageZoneCacheLocation(request.getClientIP(), ds, useDeep, track);
		}

		List<Cache>caches = selectCachesByCZ(ds, cacheLocation, track, result);

		if (caches != null) {
			return caches;
		}

		if (ds.isCoverageZoneOnly()) {
			if (ds.getGeoRedirectUrl() != null) {
				//use the NGB redirect
				caches = enforceGeoRedirect(track, ds, request.getClientIP(), null);
			} else {
				track.setResult(ResultType.MISS);
				track.setResultDetails(ResultDetails.DS_CZ_ONLY);
			}
		} else if (track.continueGeo) { 
			// continue Geo can be disabled when backup group is used -- ended up an empty cache list if reach here
			caches = selectCachesByGeo(request.getClientIP(), ds, cacheLocation, track);
		}

		return caches;
	}

	public List<Cache> selectCachesByGeo(final String clientIp, final DeliveryService deliveryService, final CacheLocation cacheLocation, final Track track) throws GeolocationException {
		Geolocation clientLocation = null;

		try {
			clientLocation = getClientLocation(clientIp, deliveryService, cacheLocation, track);
		} catch (GeolocationException e) {
			LOGGER.warn("Failed looking up Client GeoLocation: " + e.getMessage());
		}

		if (clientLocation == null) {
			if (deliveryService.getGeoRedirectUrl() != null) {
				//will use the NGB redirect
				LOGGER.debug(String
						.format("client is blocked by geolimit, use the NGB redirect url: %s",
							deliveryService.getGeoRedirectUrl()));
				return enforceGeoRedirect(track, deliveryService, clientIp, track.getClientGeolocation());
			} else {
				track.setResultDetails(ResultDetails.DS_CLIENT_GEO_UNSUPPORTED);
				return null;
			}
		}

		if (clientLocation.isDefaultLocation() && defaultGeolocationsOverride.containsKey(clientLocation.getCountryCode())) {
			clientLocation = defaultGeolocationsOverride.get(clientLocation.getCountryCode());
		}

		final List<Cache> caches = getCachesByGeo(deliveryService, clientLocation, track);

		if (caches == null || caches.isEmpty()) {
			track.setResultDetails(ResultDetails.GEO_NO_CACHE_FOUND);
		}

		track.setResult(ResultType.GEO);
		return caches;
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public DNSRouteResult route(final DNSRequest request, final Track track) throws GeolocationException {
		track.setRouteType(RouteType.DNS, request.getHostname());

		final DeliveryService ds = selectDeliveryService(request, false);

		if (ds == null) {
			track.setResult(ResultType.STATIC_ROUTE);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		if (!ds.getRoutingName().equalsIgnoreCase(request.getHostname().split("\\.")[0])) {
			// request matched the Delivery Service but is using the wrong routing name
			track.setResult(ResultType.STATIC_ROUTE);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		final DNSRouteResult result = new DNSRouteResult();

		if (!ds.isAvailable()) {
			result.setAddresses(ds.getFailureDnsResponse(request, track));
			return result;
		}

		final CacheLocation cacheLocation = getCoverageZoneCacheLocation(request.getClientIP(), ds, false, track);
		List<Cache> caches = selectCachesByCZ(ds, cacheLocation, track);

		if (caches != null) {
			track.setResult(ResultType.CZ);
			track.setClientGeolocation(cacheLocation.getGeolocation());
			result.setAddresses(inetRecordsFromCaches(ds, caches, request));
			return result;
		}

		if (ds.isCoverageZoneOnly()) {
			track.setResult(ResultType.MISS);
			track.setResultDetails(ResultDetails.DS_CZ_ONLY);
			result.setAddresses(ds.getFailureDnsResponse(request, track));
			return result;
		}

		try {
			final List<InetRecord> inetRecords = federationRegistry.findInetRecords(ds.getId(), CidrAddress.fromString(request.getClientIP()));

			if (inetRecords != null && !inetRecords.isEmpty()) {
				result.setAddresses(inetRecords);
				track.setResult(ResultType.FED);
				return result;
			}
		} catch (NetworkNodeException e) {
			LOGGER.error("Bad client address: '" + request.getClientIP() + "'");
		}

		if (track.continueGeo) {
			caches = selectCachesByGeo(request.getClientIP(), ds, cacheLocation, track);
		}

		if (caches != null) {
			track.setResult(ResultType.GEO);
			result.setAddresses(inetRecordsFromCaches(ds, caches, request));
		} else {
			track.setResult(ResultType.MISS);
			result.setAddresses(ds.getFailureDnsResponse(request, track));
		}

		return result;
	}

	public List<InetRecord> inetRecordsFromCaches(final DeliveryService ds, final List<Cache> caches, final Request request) {
		final List<InetRecord> addresses = new ArrayList<InetRecord>();
		final int maxDnsIps = ds.getMaxDnsIps();
		List<Cache> selectedCaches;

		if (maxDnsIps > 0 && isConsistentDNSRouting()) { // only consistent hash if we must
			selectedCaches = (List<Cache>) consistentHasher.selectHashables(caches, ds.getDispersion(), request.getHostname());
		} else if (maxDnsIps > 0) {
			/*
			 * We also shuffle in NameServer when adding Records to the Message prior
			 * to sending it out, as the Records are sorted later when we fill the
			 * dynamic zone if DNSSEC is enabled. We shuffle here prior to pruning
			 * for maxDnsIps so that we ensure we are spreading load across all caches
			 * assigned to this delivery service.
			*/
			Collections.shuffle(caches, random);

			selectedCaches = new ArrayList<Cache>();

			for (final Cache cache : caches) {
				selectedCaches.add(cache);

				if (selectedCaches.size() >= maxDnsIps) {
					break;
				}
			}
		} else {
			selectedCaches = caches;
		}

		for (final Cache cache : selectedCaches) {
			addresses.addAll(cache.getIpAddresses(ds.getTtls(), zoneManager, ds.isIp6RoutingEnabled()));
		}

		return addresses;
	}

	public Geolocation getClientGeolocation(final String clientIp, final Track track, final DeliveryService deliveryService) throws GeolocationException {
		if (track.isClientGeolocationQueried()) {
			return track.getClientGeolocation();
		}

		final Geolocation clientGeolocation = getLocation(clientIp, deliveryService);
		track.setClientGeolocation(clientGeolocation);
		track.setClientGeolocationQueried(true);

		return clientGeolocation;
	}

	public Geolocation getClientLocation(final String clientIp, final DeliveryService ds, final CacheLocation cacheLocation, final Track track) throws GeolocationException {
		if (cacheLocation != null) {
			return cacheLocation.getGeolocation();
		}

		final Geolocation clientGeolocation = getClientGeolocation(clientIp, track, ds);
		return ds.supportLocation(clientGeolocation);
	}

	public List<Cache> selectCachesByCZ(final DeliveryService ds, final CacheLocation cacheLocation) {
		return selectCachesByCZ(ds, cacheLocation, null);
	}

	public List<Cache> selectCachesByCZ(final String deliveryServiceId, final String cacheLocationId, final Track track) {
		return selectCachesByCZ(cacheRegister.getDeliveryService(deliveryServiceId), cacheRegister.getCacheLocation(cacheLocationId), track);
	}

	private List<Cache> selectCachesByCZ(final DeliveryService ds, final CacheLocation cacheLocation, final Track track) {
		return selectCachesByCZ(ds, cacheLocation, track, ResultType.CZ); // ResultType.CZ was the original default before DDC
	}

	private List<Cache> selectCachesByCZ(final DeliveryService ds, final CacheLocation cacheLocation, final Track track, final ResultType result) {
		if (cacheLocation == null || ds == null || !ds.isLocationAvailable(cacheLocation)) {
			return null;
		}

		final List<Cache> caches = selectCaches(cacheLocation, ds);

		if (caches != null && track != null) {
			track.setResult(result);
			if (track.isFromBackupCzGroup()) {
				track.setResultDetails(ResultDetails.DS_CZ_BACKUP_CG);
			}
			track.setResultLocation(cacheLocation.getGeolocation());
		}

		return caches;
	}

	public HTTPRouteResult multiRoute(final HTTPRequest request, final Track track) throws MalformedURLException, GeolocationException {
		final DeliveryService entryDeliveryService = cacheRegister.getDeliveryService(request, true);

		final List<SteeringResult> steeringResults = getSteeringResults(request, track, entryDeliveryService);

		if (steeringResults == null) {
			return null;
		}

		final HTTPRouteResult routeResult = new HTTPRouteResult(true);
		routeResult.setDeliveryService(entryDeliveryService);

		final List<SteeringResult> resultsToRemove = new ArrayList<>();

		for (final SteeringResult steeringResult : steeringResults) {
			final DeliveryService ds = steeringResult.getDeliveryService();

			final List<Cache> caches = selectCaches(request, ds, track);

			if (caches != null && !caches.isEmpty()) {
				final Cache cache = consistentHasher.selectHashable(caches, ds.getDispersion(), request.getPath());
				steeringResult.setCache(cache);
			} else {
				resultsToRemove.add(steeringResult);
			}
		}

		steeringResults.removeAll(resultsToRemove);

		geoSortSteeringResults(steeringResults, request.getClientIP(), entryDeliveryService);

		for (final SteeringResult steeringResult: steeringResults) {
			routeResult.addUrl(new URL(steeringResult.getDeliveryService().createURIString(request, steeringResult.getCache())));
		}

		if (routeResult.getUrls().isEmpty()) {
			routeResult.addUrl(entryDeliveryService.getFailureHttpResponse(request, track));
		}

		return routeResult;
	}

	@SuppressWarnings({ "PMD.CyclomaticComplexity", "PMD.NPathComplexity" })
	public HTTPRouteResult route(final HTTPRequest request, final Track track) throws MalformedURLException, GeolocationException {
		track.setRouteType(RouteType.HTTP, request.getHostname());

		if (isMultiRouteRequest(request)) {
			return multiRoute(request, track);
		} else {
			return singleRoute(request, track);
		}
	}

	@SuppressWarnings({ "PMD.CyclomaticComplexity", "PMD.NPathComplexity" })
	public HTTPRouteResult singleRoute(final HTTPRequest request, final Track track) throws MalformedURLException, GeolocationException {
		final DeliveryService deliveryService = getDeliveryService(request, track);

		if (deliveryService == null) {
			return null;
		}

		final HTTPRouteResult routeResult = new HTTPRouteResult(false);

		if (!deliveryService.isAvailable()) {
			routeResult.setUrl(deliveryService.getFailureHttpResponse(request, track));
			return routeResult;
		}

		routeResult.setDeliveryService(deliveryService);

		final List<Cache> caches = selectCaches(request, deliveryService, track);

		if (caches == null || caches.isEmpty()) {
			if (track.getResult() == ResultType.GEO_REDIRECT) {
				routeResult.setUrl(new URL(deliveryService.getGeoRedirectUrl()));
				LOGGER.debug(String.format("NGB redirect to url: %s for request: %s", deliveryService.getGeoRedirectUrl()
						, request.getRequestedUrl()));
				return routeResult;
			}

			routeResult.setUrl(deliveryService.getFailureHttpResponse(request, track));
			return routeResult;
		}

		final Cache cache = consistentHasher.selectHashable(caches, deliveryService.getDispersion(), request.getPath());

		// Enforce anonymous IP blocking if a DS has anonymous blocking enabled
		// and the feature is enabled
		if (deliveryService.isAnonymousIpEnabled() && AnonymousIp.getCurrentConfig().enabled) {
			AnonymousIp.enforce(this, request, deliveryService, cache, routeResult, track);

			if (routeResult.getResponseCode() == AnonymousIp.BLOCK_CODE) {
				return routeResult;
			}
		}

		if (deliveryService.isRegionalGeoEnabled()) {
			RegionalGeo.enforce(this, request, deliveryService, cache, routeResult, track);
			return routeResult;
		}

		final String uriString = deliveryService.createURIString(request, cache);
		routeResult.setUrl(new URL(uriString));

		return routeResult;
	}

	@SuppressWarnings({"PMD.NPathComplexity"})
	private List<SteeringResult> getSteeringResults(final HTTPRequest request, final Track track, final DeliveryService entryDeliveryService) {

		if (isTlsMismatch(request, entryDeliveryService)) {
			track.setResult(ResultType.ERROR);
			track.setResultDetails(ResultDetails.DS_TLS_MISMATCH);
			return null;
		}

		final List<SteeringResult> steeringResults = consistentHashMultiDeliveryService(entryDeliveryService, request.getPath());

		if (steeringResults == null || steeringResults.isEmpty()) {
			track.setResult(ResultType.DS_MISS);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		final List<SteeringResult> toBeRemoved = new ArrayList<>();
		for (final SteeringResult steeringResult : steeringResults) {
			final DeliveryService ds = steeringResult.getDeliveryService();
			if (isTlsMismatch(request, ds)) {
				track.setResult(ResultType.ERROR);
				track.setResultDetails(ResultDetails.DS_TLS_MISMATCH);
				return null;
			}
			if (ds.isRegionalGeoEnabled()) {
				LOGGER.error("Regional Geo Blocking is not supported with multi-route delivery services.. skipping " + entryDeliveryService.getId() + "/" + ds.getId());
				toBeRemoved.add(steeringResult);
			} else if (!ds.isAvailable()) {
				toBeRemoved.add(steeringResult);
			}

		}

		steeringResults.removeAll(toBeRemoved);
		return steeringResults.isEmpty() ? null : steeringResults;
	}

	private DeliveryService getDeliveryService(final HTTPRequest request, final Track track) {
		final String xtcSteeringOption = request.getHeaderValue(XTC_STEERING_OPTION);
		final DeliveryService deliveryService = consistentHashDeliveryService(cacheRegister.getDeliveryService(request, true), request.getPath(), xtcSteeringOption);

		if (deliveryService == null) {
			track.setResult(ResultType.DS_MISS);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		if (isTlsMismatch(request, deliveryService)) {
			track.setResult(ResultType.ERROR);
			track.setResultDetails(ResultDetails.DS_TLS_MISMATCH);
			return null;
		}

		return deliveryService;
	}

	private boolean isTlsMismatch(final HTTPRequest request, final DeliveryService deliveryService) {
		if (request.isSecure() && !deliveryService.isSslEnabled()) {
			return true;
		}

		if (!request.isSecure() && !deliveryService.isAcceptHttp()) {
			return true;
		}

		return false;
	}

	protected NetworkNode getDeepNetworkNode(final String ip) {
		try {
			return NetworkNode.getDeepInstance().getNetwork(ip);
		} catch (NetworkNodeException e) {
			LOGGER.warn(e);
		}
		return null;
	}

	protected NetworkNode getNetworkNode(final String ip) {
		try {
			return NetworkNode.getInstance().getNetwork(ip);
		} catch (NetworkNodeException e) {
			LOGGER.warn(e);
		}
		return null;
	}

	public CacheLocation getCoverageZoneCacheLocation(final String ip, final String deliveryServiceId) {
		return getCoverageZoneCacheLocation(ip, deliveryServiceId, false, null); // default is not deep
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	public CacheLocation getCoverageZoneCacheLocation(final String ip, final String deliveryServiceId, final boolean useDeep, final Track track) {
		final NetworkNode networkNode = useDeep ? getDeepNetworkNode(ip) : getNetworkNode(ip);
		final LocalizationMethod localizationMethod = useDeep ? LocalizationMethod.DEEP_CZ : LocalizationMethod.CZ;

		if (networkNode == null) {
			return null;
		}

		final DeliveryService deliveryService = cacheRegister.getDeliveryService(deliveryServiceId);
		CacheLocation cacheLocation = networkNode.getCacheLocation();

		if (useDeep && cacheLocation != null) {
			// lazily load deep Caches into the deep CacheLocation
			cacheLocation.loadDeepCaches(networkNode.getDeepCacheNames(), cacheRegister);
		}

		if (cacheLocation != null && !cacheLocation.isEnabledFor(localizationMethod)) {
			return null;
		}

		if (cacheLocation != null && !getSupportingCaches(cacheLocation.getCaches(), deliveryService).isEmpty()) {
			return cacheLocation;
		}

		if (useDeep) {
			// there were no available deep caches in the deep CZF
			return null;
		}

		if (networkNode.getLoc() == null) {
			return null;
		}

		// find CacheLocation
		cacheLocation = getCacheRegister().getCacheLocationById(networkNode.getLoc());
		if (cacheLocation != null && !cacheLocation.isEnabledFor(localizationMethod)) {
			track.continueGeo = false; // hit in the CZF but the cachegroup doesn't allow CZ-localization, don't fall back to GEO
			return null;
		}

		if (cacheLocation != null && !getSupportingCaches(cacheLocation.getCaches(), deliveryService).isEmpty()) {
			// lazy loading in case a CacheLocation has not yet been associated with this NetworkNode
			networkNode.setCacheLocation(cacheLocation);
			return cacheLocation;
		}

		if (cacheLocation != null && cacheLocation.getBackupCacheGroups() != null) {
			for (final String cacheGroup : cacheLocation.getBackupCacheGroups()) {
				final CacheLocation bkCacheLocation = getCacheRegister().getCacheLocationById(cacheGroup);
				if (bkCacheLocation != null && !bkCacheLocation.isEnabledFor(localizationMethod)) {
					continue;
				}
				if (bkCacheLocation != null && !getSupportingCaches(bkCacheLocation.getCaches(), deliveryService).isEmpty()) {
					LOGGER.debug("Got backup CZ cache group " + bkCacheLocation.getId() + " for " + ip + ", ds " + deliveryServiceId);
					if (track != null) {
						track.setFromBackupCzGroup(true);
					}
					return bkCacheLocation;
				}
			}
			// track.continueGeo
			// will become to false only when backups are configured and (primary group's) fallbackToClosedGeo is configured (non-empty list) to false
			// False signals subsequent cacheSelection routine to stop geo based selection.
			if (!cacheLocation.isUseClosestGeoLoc()) {
			    track.continueGeo = false;
			    return null;
			}
		} 

		// We had a hit in the CZF but the name does not match a known cache location.
		// Check whether the CZF entry has a geolocation and use it if so.
		List<CacheLocation> availableLocations = cacheRegister.filterAvailableLocations(deliveryServiceId);
		availableLocations = filterEnabledLocations(availableLocations, localizationMethod);
		final CacheLocation closestCacheLocation = getClosestCacheLocation(availableLocations, networkNode.getGeolocation(), cacheRegister.getDeliveryService(deliveryServiceId));
		if (closestCacheLocation != null) {
			LOGGER.debug("Got closest CZ cache group " + closestCacheLocation.getId() + " for " + ip + ", ds " + deliveryServiceId);
			if (track != null) {
				track.setFromBackupCzGroup(true);
			}
		}
		return closestCacheLocation;
	}

	public List<CacheLocation> filterEnabledLocations(final Collection<CacheLocation> locations, final LocalizationMethod localizationMethod) {
		return locations.stream()
				.filter(loc -> loc.isEnabledFor(localizationMethod))
				.collect(Collectors.toList());
	}

	public CacheLocation getDeepCoverageZoneCacheLocation(final String ip, final DeliveryService deliveryService) {
		return getCoverageZoneCacheLocation(ip, deliveryService, true, null);
	}

	protected CacheLocation getCoverageZoneCacheLocation(final String ip, final DeliveryService deliveryService, final boolean useDeep, final Track track) {
		return getCoverageZoneCacheLocation(ip, deliveryService.getId(), useDeep, track);
	}

	protected CacheLocation getCoverageZoneCacheLocation(final String ip, final DeliveryService deliveryService) {
		return getCoverageZoneCacheLocation(ip, deliveryService.getId());
	}

	public Cache consistentHashForCoverageZone(final String ip, final String deliveryServiceId, final String requestPath) {
		return consistentHashForCoverageZone(ip, deliveryServiceId, requestPath, false);
	}

	public Cache consistentHashForCoverageZone(final String ip, final String deliveryServiceId, final String requestPath, final boolean useDeep) {
		final DeliveryService deliveryService = cacheRegister.getDeliveryService(deliveryServiceId);
		if (deliveryService == null) {
			LOGGER.error("Failed getting delivery service from cache register for id '" + deliveryServiceId + "'");
			return null;
		}

		final CacheLocation coverageZoneCacheLocation = getCoverageZoneCacheLocation(ip, deliveryService, useDeep, null);
		final List<Cache> caches = selectCachesByCZ(deliveryService, coverageZoneCacheLocation);

		if (caches == null || caches.isEmpty()) {
			return null;
		}

		return consistentHasher.selectHashable(caches, deliveryService.getDispersion(), requestPath);
	}

	public Cache consistentHashForGeolocation(final String ip, final String deliveryServiceId, final String requestPath) {
		final DeliveryService deliveryService = cacheRegister.getDeliveryService(deliveryServiceId);
		if (deliveryService == null) {
			LOGGER.error("Failed getting delivery service from cache register for id '" + deliveryServiceId + "'");
			return null;
		}

		List<Cache> caches = null;
		if (deliveryService.isCoverageZoneOnly() && deliveryService.getGeoRedirectUrl() != null) {
				//use the NGB redirect
				caches = enforceGeoRedirect(StatTracker.getTrack(), deliveryService, ip, null);
		} else {
			final CacheLocation cacheLocation = getCoverageZoneCacheLocation(ip, deliveryServiceId);

			try {
				caches = selectCachesByGeo(ip, deliveryService, cacheLocation, StatTracker.getTrack());
			} catch (GeolocationException e) {
				LOGGER.warn("Failed gettting list of caches by geolocation for ip " + ip + " delivery service id '" + deliveryServiceId + "'");
			}
		}

		if (caches == null || caches.isEmpty()) {
			return null;
		}

		return consistentHasher.selectHashable(caches, deliveryService.getDispersion(), requestPath);
	}

	private boolean isSteeringDeliveryService(final DeliveryService deliveryService) {
		return deliveryService != null && steeringRegistry.has(deliveryService.getId());
	}

	private boolean isMultiRouteRequest(final HTTPRequest request) {
		final DeliveryService deliveryService = cacheRegister.getDeliveryService(request, true);

		if (deliveryService == null || !isSteeringDeliveryService(deliveryService)) {
			return false;
		}

		return steeringRegistry.get(deliveryService.getId()).isClientSteering();
	}

	protected Geolocation getClientLocationByCoverageZoneOrGeo(final String clientIP, final DeliveryService deliveryService) {
		Geolocation clientLocation;
		final NetworkNode networkNode = getNetworkNode(clientIP);
		if (networkNode != null && networkNode.getGeolocation() != null) {
			clientLocation = networkNode.getGeolocation();
		} else {
			try {
				clientLocation = getLocation(clientIP, deliveryService);
			} catch (GeolocationException e) {
				clientLocation = null;
			}
		}
		return deliveryService.supportLocation(clientLocation);
	}

	protected void geoSortSteeringResults(final List<SteeringResult> steeringResults, final String clientIP, final DeliveryService deliveryService) {
		if (clientIP == null || clientIP.isEmpty()
				|| steeringResults.stream().allMatch(t -> t.getSteeringTarget().getGeolocation() == null)) {
			return;
		}

		final Geolocation clientLocation = getClientLocationByCoverageZoneOrGeo(clientIP, deliveryService);
		if (clientLocation != null) {
			Collections.sort(steeringResults, new SteeringGeolocationComparator(clientLocation));
			Collections.sort(steeringResults, Comparator.comparingInt(s -> s.getSteeringTarget().getOrder())); // re-sort by order to preserve the ordering done by ConsistentHasher
		}
	}

	public List<SteeringResult> consistentHashMultiDeliveryService(final DeliveryService deliveryService, final String requestPath) {
		if (deliveryService == null) {
			return null;
		}

		final List<SteeringResult> steeringResults = new ArrayList<>();

		if (!isSteeringDeliveryService(deliveryService)) {
			steeringResults.add(new SteeringResult(null, deliveryService));
			return steeringResults;
		}

		final Steering steering = steeringRegistry.get(deliveryService.getId());
		final List<SteeringTarget> steeringTargets = consistentHasher.selectHashables(steering.getTargets(), requestPath);

		for (final SteeringTarget steeringTarget : steeringTargets) {
			final DeliveryService target = cacheRegister.getDeliveryService(steeringTarget.getDeliveryService());

			if (target != null) { // target might not be in CRConfig yet
				steeringResults.add(new SteeringResult(steeringTarget, target));
			}
		}

		return steeringResults;
	}

	public DeliveryService consistentHashDeliveryService(final String deliveryServiceId, final String requestPath) {
		return consistentHashDeliveryService(cacheRegister.getDeliveryService(deliveryServiceId), requestPath, "");
	}

	public DeliveryService consistentHashDeliveryService(final DeliveryService deliveryService, final String requestPath, final String xtcSteeringOption) {
		if (deliveryService == null) {
			return null;
		}

		if (!isSteeringDeliveryService(deliveryService)) {
			return deliveryService;
		}

		final Steering steering = steeringRegistry.get(deliveryService.getId());

		if (xtcSteeringOption != null && !xtcSteeringOption.isEmpty()) {
			return steering.hasTarget(xtcSteeringOption) ? cacheRegister.getDeliveryService(xtcSteeringOption) : null;
		}

		final String bypassDeliveryServiceId = steering.getBypassDestination(requestPath);
		if (bypassDeliveryServiceId != null && !bypassDeliveryServiceId.isEmpty()) {
			final DeliveryService bypass = cacheRegister.getDeliveryService(bypassDeliveryServiceId);
			if (bypass != null) { // bypass DS target might not be in CRConfig yet. Until then, try existing targets
				return bypass;
			}
		}

		// only select from targets in CRConfig
		final List<SteeringTarget> availableTargets = steering.getTargets().stream()
				.filter(target -> cacheRegister.getDeliveryService(target.getDeliveryService()) != null)
				.collect(Collectors.toList());
		final SteeringTarget steeringTarget = consistentHasher.selectHashable(availableTargets, deliveryService.getDispersion(), requestPath);
		return cacheRegister.getDeliveryService(steeringTarget.getDeliveryService());
	}

	/**
	 * Returns a list {@link CacheLocation}s sorted by distance from the client.
	 * If the client's location could not be determined, then the list is
	 * unsorted.
	 * 
	 * @param cacheLocations
	 *            the collection of CacheLocations to order
	 * @return the ordered list of locations
	 */
	public List<CacheLocation> orderCacheLocations(final List<CacheLocation> cacheLocations, final Geolocation clientLocation) {
		Collections.sort(cacheLocations, new CacheLocationComparator(clientLocation));
		return cacheLocations;
	}

	private CacheLocation getClosestCacheLocation(final List<CacheLocation> cacheLocations, final Geolocation clientLocation, final DeliveryService deliveryService) {
		if (clientLocation == null) {
			return null;
		}

		final List<CacheLocation> orderedLocations = orderCacheLocations(cacheLocations, clientLocation);

		for (final CacheLocation cacheLocation : orderedLocations) {
			if (!getSupportingCaches(cacheLocation.getCaches(), deliveryService).isEmpty()) {
				return cacheLocation;
			}
		}

		return null;
	}

	/*
	 * Selects a {@link Cache} from the {@link CacheLocation} provided.
	 * 
	 * @param location
	 *            the caches that will considered
	 * @param ds
	 *            the delivery service for the request
	 * @param request
	 *            the request to consider for cache selection
	 * @return the selected cache or null if none can be found
	 */
	private List<Cache> selectCaches(final CacheLocation location, final DeliveryService ds) {
		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug("Trying location: " + location.getId());
		}

		final List<Cache> caches = getSupportingCaches(location.getCaches(), ds);
		if (caches.isEmpty()) {
			if (LOGGER.isDebugEnabled()) {
				LOGGER.debug("No online, supporting caches were found at location: "
						+ location.getId());
			}
			return null;
		}

		return caches;
	}

	public Zone getZone(final Name qname, final int qtype, final InetAddress clientAddress, final boolean isDnssecRequest, final DNSAccessRecord.Builder builder) {
		return zoneManager.getZone(qname, qtype, clientAddress, isDnssecRequest, builder);
	}

	public void setRequestHeaders(final Set<String> requestHeaders) {
		this.requestHeaders = requestHeaders;
	}

	public Set<String> getRequestHeaders() {
		return requestHeaders;
	}

	public boolean isConsistentDNSRouting() {
		return consistentDNSRouting;
	}

	private List<Cache> enforceGeoRedirect(final Track track, final DeliveryService ds, final String clientIp, final Geolocation queriedClientLocation) {
		final String urlType = ds.getGeoRedirectUrlType();
		track.setResult(ResultType.GEO_REDIRECT);

		if ("NOT_DS_URL".equals(urlType)) {
			// redirect url not belongs to this DS, just redirect it
			return null;
		}

		if (!"DS_URL".equals(urlType)) {
			LOGGER.error("invalid geo redirect url type '" + urlType + "'");
			track.setResult(ResultType.MISS);
			track.setResultDetails(ResultDetails.GEO_NO_CACHE_FOUND);
			return null;
		}

		Geolocation clientLocation = queriedClientLocation;

		//redirect url belongs to this DS, will try return the caches
		if (clientLocation == null) {
			try {
				clientLocation = getLocation(clientIp, ds);
			} catch (GeolocationException e) {
				LOGGER.warn("Failed getting geolocation for client ip " + clientIp + " and delivery service '" + ds.getId() + "'");
			}
		}

		if (clientLocation == null) {
			clientLocation = ds.getMissLocation();
		}

		if (clientLocation == null) {
			LOGGER.error("cannot find a geo location for the client: " + clientIp);
			// particular error was logged in ds.supportLocation
			track.setResult(ResultType.MISS);
			track.setResultDetails(ResultDetails.DS_CLIENT_GEO_UNSUPPORTED);
			return null;
		}

		List<Cache> caches = null;

		try {
			caches = getCachesByGeo(ds, clientLocation, track);
		} catch (GeolocationException e) {
			LOGGER.error("Failed getting caches by geolocation " + e.getMessage());
		}

		if (caches == null) {
			LOGGER.warn(String.format("No Cache found by Geo in NGB redirect"));
			track.setResult(ResultType.MISS);
			track.setResultDetails(ResultDetails.GEO_NO_CACHE_FOUND);
		}

		return caches;
	}

	public void setApplicationContext(final ApplicationContext applicationContext) throws BeansException {
		this.applicationContext = applicationContext;
	}

	public void configurationChanged() {
		if (applicationContext == null) {
			LOGGER.warn("Application Context not yet ready, skipping calling listeners of configuration change");
			return;
		}

		final Map<String, ConfigurationListener> configurationListenerMap = applicationContext.getBeansOfType(ConfigurationListener.class);
		for (final ConfigurationListener configurationListener : configurationListenerMap.values()) {
			configurationListener.configurationChanged();
		}
	}

	public void setSteeringRegistry(final SteeringRegistry steeringRegistry) {
		this.steeringRegistry = steeringRegistry;
	}
}
