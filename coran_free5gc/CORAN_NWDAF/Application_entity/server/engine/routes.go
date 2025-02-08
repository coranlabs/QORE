/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the OAI Public License, Version 1.1  (the "License"); you may not use this
 * file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 *      http://www.openairinterface.org/?page_id=698
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*
 * Author: Abdelkader Mekrache <mekrache@eurecom.fr>
 * Author: Karim Boutiba 	   <boutiba@eurecom.fr>
 * Author: Arina Prostakova    <prostako@eurecom.fr>
 * Description: This file contains server routes.
 */

package engine

import (
	"net/http"
)

// ------------------------------------------------------------------------------
// NewRouter - create router for HTTP server.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	// register routes
	mux.HandleFunc(config.Routes.NumOfUe, nwPerfNumOfUe)
	mux.HandleFunc(config.Routes.SessSuccRatio, nwPerfNumOfPdu)
	mux.HandleFunc(config.Routes.UeComm, ueComm)
	mux.HandleFunc(config.Routes.UeMob, ueMob)
	return mux
}
