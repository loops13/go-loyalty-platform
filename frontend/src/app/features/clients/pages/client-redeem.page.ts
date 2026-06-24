import { CommonModule } from '@angular/common';
import { Component, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { forkJoin } from 'rxjs';

import { ClientResp, RewardResp } from '../../../api/generated/api.models';
import { ClientsApiService } from '../../../api/generated/clients-api.service';
import { RewardsApiService } from '../../../api/generated/rewards-api.service';

@Component({
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <section class="card">
      <h2>Redeem reward</h2>
      <p class="muted" *ngIf="client()">{{ client()?.name }} · {{ client()?.pointBalance }} points</p>

      <form class="form" (ngSubmit)="submit()">
        <label>
          <span>Reward</span>
          <select name="rewardId" [ngModel]="rewardId()" (ngModelChange)="rewardId.set($event)">
            <option value="" disabled>Select a reward</option>
            <option *ngFor="let reward of rewards()" [value]="reward.id">
              {{ reward.name }} ({{ reward.pointsCost }} points)
            </option>
          </select>
        </label>

        <p class="error" *ngIf="error()">{{ error() }}</p>
        <button type="submit" [disabled]="saving()">Redeem</button>
        <a class="secondary-link" [routerLink]="['/clients', clientId]">Cancel</a>
      </form>
    </section>
  `,
})
export class ClientRedeemPage implements OnInit {
  clientId = '';
  readonly client = signal<ClientResp | null>(null);
  readonly rewards = signal<RewardResp[]>([]);
  readonly rewardId = signal('');
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);

  constructor(
    private readonly clientsApi: ClientsApiService,
    private readonly rewardsApi: RewardsApiService,
    private readonly route: ActivatedRoute,
    private readonly router: Router,
  ) {}

  ngOnInit(): void {
    this.clientId = this.route.snapshot.paramMap.get('id') ?? '';
    if (!this.clientId) {
      this.error.set('Missing client id.');
      return;
    }

    forkJoin({
      client: this.clientsApi.getClient(this.clientId),
      rewards: this.rewardsApi.listRewards(),
    }).subscribe({
      next: ({ client, rewards }) => {
        this.client.set(client);
        this.rewards.set(rewards);
        this.rewardId.set(rewards[0]?.id ?? '');
      },
      error: () => this.error.set('Unable to load redemption data.'),
    });
  }

  submit(): void {
    if (!this.clientId || !this.rewardId()) {
      this.error.set('Select a reward.');
      return;
    }

    this.saving.set(true);
    this.error.set(null);

    this.rewardsApi.redeemReward(this.clientId, { rewardId: this.rewardId() }).subscribe({
      next: () => {
        this.saving.set(false);
        this.router.navigate(['/clients', this.clientId]);
      },
      error: () => {
        this.saving.set(false);
        this.error.set('Unable to redeem reward.');
      },
    });
  }
}
