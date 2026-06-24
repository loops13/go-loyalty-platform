import { CommonModule } from '@angular/common';
import { Component, OnInit, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { forkJoin } from 'rxjs';

import { ClientResp, AwardResp } from '../../../api/generated/api.models';
import { ClientsApiService } from '../../../api/generated/clients-api.service';

@Component({
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <section class="card" *ngIf="client(); else loadingTpl">
      <div class="header">
        <div>
          <h2>{{ client()?.name }}</h2>
          <p class="muted">{{ client()?.email }}</p>
        </div>
        <p class="balance">{{ client()?.pointBalance }} points</p>
      </div>

      <div class="actions">
        <a [routerLink]="['/clients', clientId, 'award']">Award points</a>
        <a [routerLink]="['/clients', clientId, 'redeem']">Redeem reward</a>
      </div>

      <h3>Transaction history</h3>
      <p class="muted" *ngIf="!awards().length">No awards yet.</p>
      <ul class="list">
        <li *ngFor="let award of awards()">
          <strong>{{ award.type }}</strong>
          <span>{{ award.pointsAwarded }} points</span>
          <small>{{ award.createdAt | date:'medium' }}</small>
        </li>
      </ul>
    </section>

    <ng-template #loadingTpl>
      <section class="card">
        <p class="muted" *ngIf="loading()">Loading...</p>
        <p class="error" *ngIf="error()">{{ error() }}</p>
      </section>
    </ng-template>
  `,
})
export class ClientDetailPage implements OnInit {
  clientId = '';
  readonly client = signal<ClientResp | null>(null);
  readonly awards = signal<AwardResp[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  constructor(private readonly api: ClientsApiService, private readonly route: ActivatedRoute) {}

  ngOnInit(): void {
    this.clientId = this.route.snapshot.paramMap.get('id') ?? '';
    this.load();
  }

  load(): void {
    if (!this.clientId) {
      this.error.set('Missing client id.');
      return;
    }

    this.loading.set(true);
    this.error.set(null);

    forkJoin({
      client: this.api.getClient(this.clientId),
      awards: this.api.getClientAwards(this.clientId),
    }).subscribe({
      next: ({ client, awards }) => {
        this.client.set(client);
        this.awards.set(awards);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Unable to load client details.');
        this.loading.set(false);
      },
    });
  }
}
