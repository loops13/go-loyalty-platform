import { CommonModule } from '@angular/common';
import { Component, OnInit, signal } from '@angular/core';

import { RewardResp } from '../../../api/generated/api.models';
import { RewardsApiService } from '../../../api/generated/rewards-api.service';

@Component({
  standalone: true,
  imports: [CommonModule],
  template: `
    <section class="card">
      <div class="header">
        <div>
          <h2>Rewards</h2>
          <p class="muted">Available rewards that clients can redeem.</p>
        </div>
        <button type="button" (click)="load()">Refresh</button>
      </div>

      <p class="muted" *ngIf="loading()">Loading...</p>
      <p class="error" *ngIf="error()">{{ error() }}</p>

      <ul class="list">
        <li *ngFor="let reward of rewards()">
          <strong>{{ reward.name }}</strong>
          <span>{{ reward.pointsCost }} points</span>
        </li>
      </ul>
    </section>
  `,
})
export class RewardListPage implements OnInit {
  readonly rewards = signal<RewardResp[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  constructor(private readonly api: RewardsApiService) {}

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.loading.set(true);
    this.error.set(null);

    this.api.listRewards().subscribe({
      next: (rewards) => {
        this.rewards.set(rewards);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Unable to load rewards.');
        this.loading.set(false);
      },
    });
  }
}
