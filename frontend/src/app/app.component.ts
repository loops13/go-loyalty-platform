import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive, RouterOutlet],
  template: `
    <div class="shell">
      <header class="topbar">
        <div>
          <p class="eyebrow">Go Loyalty Platform</p>
          <h1>Rewards dashboard</h1>
        </div>
        <nav class="nav">
          <a routerLink="/clients" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">Clients</a>
          <a routerLink="/clients/new" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">Create client</a>
          <a routerLink="/rewards" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">Rewards</a>
        </nav>
      </header>

      <main class="content">
        <router-outlet></router-outlet>
      </main>
    </div>
  `,
  styles: [`
    .shell { min-height: 100vh; padding: 24px; background: #f7f7fb; color: #111827; }
    .topbar { display: flex; justify-content: space-between; align-items: center; gap: 16px; max-width: 1120px; margin: 0 auto 24px; }
    .eyebrow { margin: 0 0 4px; font-size: 12px; text-transform: uppercase; letter-spacing: .12em; color: #6b7280; }
    h1 { margin: 0; }
    .nav { display: flex; gap: 12px; flex-wrap: wrap; }
    .nav a { text-decoration: none; padding: 8px 12px; border-radius: 999px; color: #374151; background: #fff; border: 1px solid #e5e7eb; }
    .nav a.active { background: #2563eb; color: #fff; border-color: #2563eb; }
    .content { max-width: 1120px; margin: 0 auto; }
  `],
})
export class AppComponent {}
