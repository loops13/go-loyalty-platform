import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { forkJoin } from 'rxjs';

import { ApiService } from './api.service';
import { Client, Reward } from './models';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  clients: Client[] = [];
  rewards: Reward[] = [];
  loading = false;
  error = '';
  name = '';
  email = '';

  constructor(private readonly api: ApiService) {}

  ngOnInit(): void {
    this.refresh();
  }

  refresh(): void {
    this.loading = true;
    this.error = '';

    forkJoin({
      clients: this.api.listClients(),
      rewards: this.api.listRewards()
    }).subscribe({
      next: ({ clients, rewards }) => {
        this.clients = clients;
        this.rewards = rewards;
        this.loading = false;
      },
      error: () => {
        this.error = 'Unable to load data from the backend.';
        this.loading = false;
      }
    });
  }

  createClient(): void {
    const name = this.name.trim();
    const email = this.email.trim();

    if (!name || !email) {
      this.error = 'Name and email are required.';
      return;
    }

    this.error = '';
    this.api.createClient({ name, email }).subscribe({
      next: () => {
        this.name = '';
        this.email = '';
        this.refresh();
      },
      error: () => {
        this.error = 'Unable to create client.';
      }
    });
  }
}
