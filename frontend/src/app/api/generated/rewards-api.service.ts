import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { RedeemReq, RedeemResp, RewardResp } from './api.models';
import { ApiBaseService } from './api-base.service';

@Injectable({ providedIn: 'root' })
export class RewardsApiService extends ApiBaseService {
  listRewards(): Observable<RewardResp[]> {
    return this.http.get<RewardResp[]>(this.url('/rewards'));
  }

  redeemReward(clientId: string, payload: RedeemReq): Observable<RedeemResp> {
    return this.http.post<RedeemResp>(this.url(`/clients/${clientId}/redeem`), payload);
  }
}
