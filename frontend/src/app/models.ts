export interface Client {
  id: string;
  name: string;
  email: string;
  pointBalance: number;
}

export interface Reward {
  id: string;
  name: string;
  pointsCost: number;
}

export interface CreateClientRequest {
  name: string;
  email: string;
}
