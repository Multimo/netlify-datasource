import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export interface NetlifyQuery extends DataQuery {
  siteId?: string;
  entity?: string;
  parsingOptions?: {
    selectedFields: string[]
  }
}

/**
 * These are options configured for each DataSource instance
 */
export interface NetlifyDataSourceOptions extends DataSourceJsonData {
  path?: string;
  accountId?: string;
  siteId?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface NetlifySecureJsonData {
  accessToken?: string;
}
