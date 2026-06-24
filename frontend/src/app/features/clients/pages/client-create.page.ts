import { CommonModule } from '@angular/common';
import { Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

import { ClientsApiService } from '../../../api/generated/clients-api.service';

@Component({
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <section class="card">
      <h2>Create client</h2>
      <p class="muted">Create a new client profile that can receive points and redeem rewards.</p>

      <form class="form" (ngSubmit)="submit()">
        <label>
          <span>Name</span>
          <input name="name" [ngModel]="name()" (ngModelChange)="name.set($event)" />
        </label>
        <label>
          <span>Email</span>
          <input name="email" [ngModel]="email()" (ngModelChange)="email.set($event)" />
        </label>

        <p class="error" *ngIf="error()">{{ error() }}</p>
        <button type="submit" [disabled]="saving()">Create</button>
      </form>
    </section>
  `,
})
export class ClientCreatePage {
  readonly name = signal('');
  readonly email = signal('');
  readonly saving = signal(false);
  readonly error = signal<string | null>(null);

  constructor(private readonly api: ClientsApiService, private readonly router: Router) {}

  submit(): void {
    const name = this.name().trim();
    const email = this.email().trim();

    if (!name || !email) {
      this.error.set('Name and email are required.');
      return;
    }

    this.saving.set(true);
    this.error.set(null);

    this.api.createClient({ name, email }).subscribe({
      next: (client) => {
        this.saving.set(false);
        this.router.navigate(['/clients', client.id]);
      },
      error: () => {
        this.saving.set(false);
        this.error.set('Unable to create client.');
      },
    });
  }
}
