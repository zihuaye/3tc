/*
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
import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';

import { first } from 'rxjs/operators';

import { APIService } from '../../services';
import { DeliveryService } from '../../models/deliveryservice';

@Component({
	selector: 'dash',
	templateUrl: './dashboard.component.html',
	styleUrls: ['./dashboard.component.scss']
})
/**
 * Controller for the dashboard. Doesn't do much yet.
*/
export class DashboardComponent implements OnInit {
	deliveryServices: DeliveryService[];
	loading = true;

	// Fuzzy search control
	fuzzControl = new FormControl('');

	constructor (private readonly api: APIService) { }

	ngOnInit () {
		this.api.getDeliveryServices().pipe(first()).subscribe(
			r => {
				this.deliveryServices = r;
				this.loading = false;
			}
		);
	}

	/**
	 * Checks if a Delivery Service matches a fuzzy search term
	 * @param ds The Delivery Service being checked
	 * @returns `true` if `ds` matches the fuzzy search term, `false` otherwise
	*/
	fuzzy (ds: DeliveryService): boolean {
		if (!this.fuzzControl.value) {
			return true;
		}
		const testVal = ds.displayName.toLocaleLowerCase();
		let n = -1;
		for (const l of this.fuzzControl.value.toLocaleLowerCase()) {
			/* tslint:disable */
			if (!~(n = testVal.indexOf(l, n + 1))) {
			/* tslint:enable */
				return false;
			}
		}
		return true;
	}

}
