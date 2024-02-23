import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
// import { VariableEditor } from './components/VariableEditor';
import { NetlifyQuery, NetlifyDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, NetlifyQuery, NetlifyDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
// .setVariableQueryEditor(VariableEditor)

