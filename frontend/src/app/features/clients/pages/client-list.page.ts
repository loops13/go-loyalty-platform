import { CommonModule } from '@angular/common';
import { Component, OnInit, signal } from '@angular/core';
import { RouterLink } from '@angular/router';

import { ClientsApiService } from '../../../api/generated/clients-api.service';
import { ClientResp } from '../../../api/generated/api.models';

@Component({
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <section class="card">
      <div class="header">
        <div>
          <h2>Clients</h2>
          <p class="muted">Browse all clients and open a profile for details, awards, and redemption.</p>
        </div>
        <button type="button" (click)="load()">Refresh</button>
      </div>

      <p class="muted" *ngIf="loading()">Loading...</p>
      <p class="error" *ngIf="error()">{{ error() }}</p>

      <table *ngIf="clients().length" class="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Point balance</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr *ngFor="let client of clients()">
            <td>{{ client.name }}</td>
            <td>{{ client.email }}</td>
            <td>{{ client.pointBalance }}</td>
            <td class="row-actions">
              <a [routerLink]="['/clients', client.id]">Open</a>
              <a href="javascript:void(0)" class="action-link danger" (click)="remove(client)">Delete client</a>
            </td>
          </tr>
        </tbody>
      </table>
    </section>
  `,
})
export class ClientListPage implements OnInit {
  readonly clients = signal<ClientResp[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  constructor(private readonly api: ClientsApiService) {}

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.loading.set(true);
    this.error.set(null);

    this.api.listClients().subscribe({
      next: (clients) => {
        this.clients.set(clients);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Unable to load clients.');
        this.loading.set(false);
      },
    });
  }

  remove(client: ClientResp): void {
    if (!confirm(`Delete ${client.name}? This will remove their awards too.`)) {
      return;
    }

    this.api.deleteClient(client.id).subscribe({
      next: () => this.load(),
      error: () => this.error.set('Unable to delete client.'),
    });
  }
}
