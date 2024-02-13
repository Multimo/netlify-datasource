import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { NetlifyQuery, NetlifyDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, NetlifyQuery, NetlifyDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
