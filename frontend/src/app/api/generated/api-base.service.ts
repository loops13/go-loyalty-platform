import { HttpClient } from '@angular/common/http';
import { inject } from '@angular/core';

import { API_BASE_URL } from '../../core/config/api-base-url.token';

export abstract class ApiBaseService {
  protected readonly http = inject(HttpClient);
  protected readonly baseUrl = inject(API_BASE_URL);

  protected url(path: string): string {
    return `${this.baseUrl}${path}`;
  }
}
