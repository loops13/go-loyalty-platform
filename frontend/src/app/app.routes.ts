import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'clients',
  },
  {
    path: 'clients',
    loadComponent: () =>
      import('./features/clients/pages/client-list.page').then((m) => m.ClientListPage),
  },
  {
    path: 'clients/new',
    loadComponent: () =>
      import('./features/clients/pages/client-create.page').then((m) => m.ClientCreatePage),
  },
  {
    path: 'clients/:id',
    loadComponent: () =>
      import('./features/clients/pages/client-detail.page').then((m) => m.ClientDetailPage),
  },
  {
    path: 'clients/:id/award',
    loadComponent: () =>
      import('./features/clients/pages/client-award.page').then((m) => m.ClientAwardPage),
  },
  {
    path: 'clients/:id/redeem',
    loadComponent: () =>
      import('./features/clients/pages/client-redeem.page').then((m) => m.ClientRedeemPage),
  },
  {
    path: 'rewards',
    loadComponent: () =>
      import('./features/rewards/pages/reward-list.page').then((m) => m.RewardListPage),
  },
];
