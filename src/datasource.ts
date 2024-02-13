import { DataSourceInstanceSettings, CoreApp } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { NetlifyQuery, NetlifyDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<NetlifyQuery, NetlifyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NetlifyDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<NetlifyQuery> {
    return DEFAULT_QUERY;
  }
}
