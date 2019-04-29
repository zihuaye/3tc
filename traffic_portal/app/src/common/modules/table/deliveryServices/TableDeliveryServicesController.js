/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var TableDeliveryServicesController = function(deliveryServices, $anchorScroll, $scope, $state, $location, $uibModal, $window, deliveryServiceService, deliveryServiceRequestService, dateUtils, deliveryServiceUtils, locationUtils, messageModel, propertiesModel, userModel) {

    var protocols = deliveryServiceUtils.protocols;

    var qstrings = deliveryServiceUtils.qstrings;

    var dsRequestsEnabled = propertiesModel.properties.dsRequests.enabled;

    var createDeliveryService = function(typeName) {
        var path = '/delivery-services/new?type=' + typeName;
        locationUtils.navigateToPath(path);
    };

    var clone = function(ds) {
        var params = {
            title: 'Clone Delivery Service: ' + ds.xmlId,
            message: "Please select a content routing category for the clone"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
            controller: 'DialogSelectController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function() {
                    // the following represent the 4 categories of delivery services
                    // the ids are arbitrary but the dialog.select dropdown needs them
                    return [
                        { id: 1, name: 'ANY_MAP' },
                        { id: 2, name: 'DNS' },
                        { id: 3, name: 'HTTP' },
                        { id: 4, name: 'STEERING' }
                    ];
                }
            }
        });
        modalInstance.result.then(function(type) {
            locationUtils.navigateToPath('/delivery-services/' + ds.id + '/clone?type=' + type.name);
        });
    };

    var confirmDelete = function(deliveryService) {
        var params = {
            title: 'Delete Delivery Service: ' + deliveryService.xmlId,
            key: deliveryService.xmlId
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/delete/dialog.delete.tpl.html',
            controller: 'DialogDeleteController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            if (dsRequestsEnabled) {
                createDeliveryServiceDeleteRequest(deliveryService);
            } else {
                deliveryServiceService.deleteDeliveryService(deliveryService)
                    .then(
                        function() {
                            messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + deliveryService.xmlId + ' ] deleted' } ], false);
                            $scope.refresh();                        },
                        function(fault) {
                            $anchorScroll(); // scrolls window to top
                            messageModel.setMessages(fault.data.alerts, false);
                        }
                    );
            }
        }, function () {
            // do nothing
        });
    };

    var createDeliveryServiceDeleteRequest = function(deliveryService) {
        var params = {
            title: "Delivery Service Delete Request",
            message: 'All delivery service deletions must be reviewed.'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/deliveryServiceRequest/dialog.deliveryServiceRequest.tpl.html',
            controller: 'DialogDeliveryServiceRequestController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                statuses: function() {
                    var statuses = [
                        { id: $scope.DRAFT, name: 'Save Request as Draft' },
                        { id: $scope.SUBMITTED, name: 'Submit Request for Review and Deployment' }
                    ];
                    if (userModel.user.roleName == propertiesModel.properties.dsRequests.roleNeededToSkip) {
                        statuses.push({ id: $scope.COMPLETE, name: 'Fulfill Request Immediately' });
                    }
                    return statuses;
                }
            }
        });
        modalInstance.result.then(function(options) {
            var status = 'draft';
            if (options.status.id == $scope.SUBMITTED || options.status.id == $scope.COMPLETE) {
                status = 'submitted';
            };

            var dsRequest = {
                changeType: 'delete',
                status: status,
                deliveryService: deliveryService
            };

            // if the user chooses to complete/fulfill the delete request immediately, the ds will be deleted and behind the
            // scenes a delivery service request will be created and marked as complete
            if (options.status.id == $scope.COMPLETE) {
                // first delete the ds
                deliveryServiceService.deleteDeliveryService(deliveryService)
                    .then(
                        function() {
                            // then create the ds request
                            deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest).
                            then(
                                function(response) {
                                    var comment = {
                                        deliveryServiceRequestId: response.id,
                                        value: options.comment
                                    };
                                    // then create the ds request comment
                                    deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).
                                    then(
                                        function() {
                                            var promises = [];
                                            // assign the ds request
                                            promises.push(deliveryServiceRequestService.assignDeliveryServiceRequest(response.id, userModel.user.id));
                                            // set the status to 'complete'
                                            promises.push(deliveryServiceRequestService.updateDeliveryServiceRequestStatus(response.id, 'complete'));
                                            // and finally refresh the delivery services table
                                            messageModel.setMessages([ { level: 'success', text: 'Delivery service [ ' + deliveryService.xmlId + ' ] deleted' } ], false);
                                            $scope.refresh();
                                        }
                                    );
                                }
                            );
                        },
                        function(fault) {
                            $anchorScroll(); // scrolls window to top
                            messageModel.setMessages(fault.data.alerts, false);
                        }
                    );
            } else {
                deliveryServiceRequestService.createDeliveryServiceRequest(dsRequest).
                    then(
                        function(response) {
                            var comment = {
                                deliveryServiceRequestId: response.id,
                                value: options.comment
                            };
                            deliveryServiceRequestService.createDeliveryServiceRequestComment(comment).
                                then(
                                    function() {
                                        messageModel.setMessages([ { level: 'success', text: 'Created request to ' + dsRequest.changeType + ' the ' + dsRequest.deliveryService.xmlId + ' delivery service' } ], true);
                                        locationUtils.navigateToPath('/delivery-service-requests');
                                    }
                                );
                        }
                    );
            }
        });
    };

    $scope.deliveryServices = deliveryServices;

    $scope.showChartsButton = propertiesModel.properties.deliveryServices.charts.customLink.show;

    $scope.openCharts = deliveryServiceUtils.openCharts;

    $scope.getRelativeTime = dateUtils.getRelativeTime;

    $scope.navigateToPath = locationUtils.navigateToPath;

    $scope.DRAFT = 0;
    $scope.SUBMITTED = 1;
    $scope.REJECTED = 2;
    $scope.PENDING = 3;
    $scope.COMPLETE = 4;

    $scope.contextMenuItems = [
        {
            text: 'Open in New Tab',
            click: function ($itemScope) {
                $window.open('/#!/delivery-services/' + $itemScope.ds.id + '?type=' + $itemScope.ds.type, '_blank');
            }
        },
        null, // Divider
        {
            text: 'Edit',
            click: function ($itemScope) {
                $scope.editDeliveryService($itemScope.ds);
            }
        },
        {
            text: 'Clone',
            click: function ($itemScope) {
                clone($itemScope.ds);
            }
        },
        {
            text: 'Delete',
            click: function ($itemScope) {
                confirmDelete($itemScope.ds);
            }
        },
        null, // Divider
        {
            text: 'View Charts',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/charts?type=' + $itemScope.ds.type);
            }
        },
        null, // Divider
        {
            text: 'Manage SSL Keys',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/ssl-keys?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage URL Sig Keys',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/url-sig-keys?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage URI Signing Keys',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/uri-signing-keys?type=' + $itemScope.ds.type);
            }
        },
        null, // Divider
        {
            text: 'Manage Targets',
            displayed: function ($itemScope) {
                // only show for steering* delivery services
                return $itemScope.ds.type.indexOf('STEERING') != -1;
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/targets?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Origins',
            displayed: function ($itemScope) {
                // only show for non-steering* delivery services
                return $itemScope.ds.type.indexOf('STEERING') == -1;
            },
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/origins?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Servers',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/servers?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Regexes',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/regexes?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Invalidation Requests',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/jobs?type=' + $itemScope.ds.type);
            }
        },
        {
            text: 'Manage Static DNS Entries',
            click: function ($itemScope) {
                locationUtils.navigateToPath('/delivery-services/' + $itemScope.ds.id + '/static-dns-entries?type=' + $itemScope.ds.type);
            }
        }
    ];

    $scope.editDeliveryService = function(ds) {
        var path = '/delivery-services/' + ds.id + '?type=' + ds.type;
        locationUtils.navigateToPath(path);
    };

    $scope.refresh = function() {
        $state.reload(); // reloads all the resolves for the view
    };

    $scope.protocol = function(ds) {
        return protocols[ds.protocol];
    };

    $scope.qstring = function(ds) {
        return qstrings[ds.qstringIgnore];
    };

    $scope.selectDSType = function() {
        var params = {
            title: 'Create Delivery Service',
            message: "Please select a content routing category"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/select/dialog.select.tpl.html',
            controller: 'DialogSelectController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function() {
                    // the following represent the 4 categories of delivery services
                    // the ids are arbitrary but the dialog.select dropdown needs them
                    return [
                        { id: 1, name: 'ANY_MAP' },
                        { id: 2, name: 'DNS' },
                        { id: 3, name: 'HTTP' },
                        { id: 4, name: 'STEERING' }
                    ];
                }
            }
        });
        modalInstance.result.then(function(type) {
            createDeliveryService(type.name);
        }, function () {
            // do nothing
        });
    };

    $scope.compareDSs = function() {
        var params = {
            title: 'Compare Delivery Services',
            message: "Please select 2 delivery services to compare",
            label: "xmlId"
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/compare/dialog.compare.tpl.html',
            controller: 'DialogCompareController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                },
                collection: function(deliveryServiceService) {
                    return deliveryServiceService.getDeliveryServices();
                }
            }
        });
        modalInstance.result.then(function(dss) {
            $location.path($location.path() + '/compare/' + dss[0].id + '/' + dss[1].id);
        }, function () {
            // do nothing
        });
    };

    angular.element(document).ready(function () {
        $('#deliveryServicesTable').dataTable({
            "aLengthMenu": [[25, 50, 100, -1], [25, 50, 100, "All"]],
            "iDisplayLength": 25,
            "columnDefs": [
                { 'orderable': false, 'targets': 12 }
            ],
            "aaSorting": []
        });
    });

};

TableDeliveryServicesController.$inject = ['deliveryServices', '$anchorScroll', '$scope', '$state', '$location', '$uibModal', '$window', 'deliveryServiceService', 'deliveryServiceRequestService', 'dateUtils', 'deliveryServiceUtils', 'locationUtils', 'messageModel', 'propertiesModel', 'userModel'];
module.exports = TableDeliveryServicesController;
