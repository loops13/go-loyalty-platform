import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { environment } from '../environments/environment';
import { Client, CreateClientRequest, Reward } from './models';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly baseUrl = environment.apiBaseUrl;

  constructor(private readonly http: HttpClient) {}

  listClients() {
    return this.http.get<Client[]>(`${this.baseUrl}/clients`);
  }

  listRewards() {
    return this.http.get<Reward[]>(`${this.baseUrl}/rewards`);
  }

  createClient(payload: CreateClientRequest) {
    return this.http.post<Client>(`${this.baseUrl}/clients`, payload);
  }
}
