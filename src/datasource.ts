import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { NetlifyQuery, NetlifyDataSourceOptions } from './types';
import { VariableSupport } from 'variables/VariableSupport';


export class DataSource extends DataSourceWithBackend<NetlifyQuery, NetlifyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<NetlifyDataSourceOptions>) {
    super(instanceSettings);
    this.annotations = {};
    this.variables = new VariableSupport(this);
  }

  async getSiteIds(): Promise<string[]> {
    return await this.getResource('sites');
  }

  applyTemplateVariables(query: NetlifyQuery, scopedVars: ScopedVars): NetlifyQuery {
    const templateSrv = getTemplateSrv();
    const siteIds = templateSrv.replace(query.siteId, scopedVars);
    console.log('applyTemplateVariables', { query, scopedVars, siteIds })
    return {
      ...query,
      siteId: scopedVars.siteId?.value ?? siteIds ?? query.siteId,
    };
  }
}
