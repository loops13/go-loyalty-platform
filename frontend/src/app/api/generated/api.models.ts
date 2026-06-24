export type AwardType =
  | 'MONTHLY_CONTRIBUTION'
  | 'INCREASE_CONTRIBUTION'
  | 'LOYALTY_12_MONTHS'
  | 'FINANCIAL_ASSESSMENT';

export interface ClientResp {
  id: string;
  name: string;
  email: string;
  pointBalance: number;
}

export interface CreateClientReq {
  name: string;
  email: string;
}

export interface AwardReq {
  type: AwardType;
}

export interface AwardResp {
  id: string;
  clientId: string;
  type: string;
  pointsAwarded: number;
  createdAt: string;
}

export interface RewardResp {
  id: string;
  name: string;
  pointsCost: number;
}

export interface RedeemReq {
  rewardId: string;
}

export interface RedeemResp {
  reward: RewardResp;
  balance: number;
}

export interface ErrorResp {
  code: string;
  message: string;
}

export const AWARD_TYPE_OPTIONS: Array<{ value: AwardType; label: string }> = [
  { value: 'MONTHLY_CONTRIBUTION', label: 'Monthly investment contribution' },
  { value: 'INCREASE_CONTRIBUTION', label: 'Increase contribution amount' },
  { value: 'LOYALTY_12_MONTHS', label: 'Keep investment active for 12 months' },
  { value: 'FINANCIAL_ASSESSMENT', label: 'Complete financial wellness assessment' },
];
