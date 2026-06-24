import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import {
  AwardReq,
  AwardResp,
  ClientResp,
  CreateClientReq,
} from './api.models';
import { ApiBaseService } from './api-base.service';

@Injectable({ providedIn: 'root' })
export class ClientsApiService extends ApiBaseService {
  listClients(): Observable<ClientResp[]> {
    return this.http.get<ClientResp[]>(this.url('/clients'));
  }

  createClient(payload: CreateClientReq): Observable<ClientResp> {
    return this.http.post<ClientResp>(this.url('/clients'), payload);
  }

  getClient(clientId: string): Observable<ClientResp> {
    return this.http.get<ClientResp>(this.url(`/clients/${clientId}`));
  }

  getClientAwards(clientId: string): Observable<AwardResp[]> {
    return this.http.get<AwardResp[]>(this.url(`/clients/${clientId}/awards`));
  }

  awardPoints(clientId: string, payload: AwardReq): Observable<AwardResp> {
    return this.http.post<AwardResp>(this.url(`/clients/${clientId}/awards`), payload);
  }
}
