import { CommonModule } from '@angular/common';
import { Component, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';

import {
  AWARD_TYPE_OPTIONS,
  AwardType,
  ClientResp,
} from '../../../api/generated/api.models';
import { ClientsApiService } from '../../../api/generated/clients-api.service';

@Component({
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <section class="card">
      <h2>Award points</h2>
      <p class="muted" *ngIf="client()">{{ client()?.name }} · {{ client()?.email }}</p>

      <form class="form" (ngSubmit)="submit()">
        <label>
          <span>Award type</span>
          <select name="type" [ngModel]="awardType()" (ngModelChange)="awardType.set($event)">
            <option *ngFor="let option of awardOptions" [value]="option.value">{{ option.label }}</option>
          </select>
        </label>

        <p class="error" *ngIf="error()">{{ error() }}</p>
        <button type="submit" [disabled]="saving()">Award points</button>
        <a class="secondary-link" [routerLink]="['/clients', clientId]">Cancel</a>
      </form>
    </section>
  `,
})
export class ClientAwardPage implements OnInit {
  readonly awardOptions = AWARD_TYPE_OPTIONS;
  clientId = '';
  readonly client = signal<ClientResp | null>(null);
  readonly awardType = signal<AwardType>('MONTHLY_CONTRIBUTION');
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);

  constructor(private readonly api: ClientsApiService, private readonly route: ActivatedRoute, private readonly router: Router) {}

  ngOnInit(): void {
    this.clientId = this.route.snapshot.paramMap.get('id') ?? '';
    if (!this.clientId) {
      this.error.set('Missing client id.');
      return;
    }

    this.api.getClient(this.clientId).subscribe({
      next: (client) => this.client.set(client),
      error: () => this.error.set('Unable to load client.'),
    });
  }

  submit(): void {
    if (!this.clientId) {
      this.error.set('Missing client id.');
      return;
    }

    this.saving.set(true);
    this.error.set(null);

    this.api.awardPoints(this.clientId, { type: this.awardType() }).subscribe({
      next: () => {
        this.saving.set(false);
        this.router.navigate(['/clients', this.clientId]);
      },
      error: () => {
        this.saving.set(false);
        this.error.set('Unable to award points.');
      },
    });
  }
}
