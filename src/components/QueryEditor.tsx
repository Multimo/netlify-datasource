import React, { useCallback, useEffect, useState } from 'react';
import { HorizontalGroup, InlineField, Input, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource';
import { NetlifyDataSourceOptions, NetlifyQuery } from '../types';
import { ParsingOptionsEditor } from './TransfromEditor';
import { ParametersEditor } from './ParametersEditor';

type Props = QueryEditorProps<DataSource, NetlifyQuery, NetlifyDataSourceOptions>;

const entity_options = [
  { label: 'Builds', value: 'builds', description: 'Query for list of all builds by site id' },
  { label: 'Build Account', value: 'builds-account', description: 'Query Current Build Status by Account id' },
  { label: 'Deployments', value: 'deployments', description: 'Query for list of all deployments by site id' },
  { label: 'Forms', value: 'forms', description: 'Query for list of all forms by site id' },
  { label: 'Form Submissions', value: 'form-submissions', description: 'Query for list of form submissions by site id' },
  { label: 'Sites', value: 'sites', description: 'Query for list of owned Sites' },
  { label: 'Accounts', value: 'accounts', description: 'Query for list of Accounts' },
];

const entities_requiring_site_id = ['builds', 'deployments', 'forms', 'form-submissions']

const default_site_id = { label: 'Default Site Id', value: '' }

export function QueryEditor({ query, onChange, onRunQuery, datasource, data, ...rest }: Props) {
  console.log({ query, onChange, onRunQuery, datasource, ...rest })
  const [_, setSiteIdOptions] = useState([default_site_id])

  const handleEntityChange = (value: SelectableValue<string>) => {
    onChange({ ...query, entity: value.value });
    onRunQuery();
  }

  const updateSiteIds = useCallback(() => {
    return datasource.getSiteIds().then((siteIds) => {
      console.log({ siteIds })
      const options = siteIds.map((site) => ({ label: site, value: site }))
      setSiteIdOptions(options)
    }).catch(console.error)
  }, [datasource])


  const handleSiteIdUpdates = (value: string) => {
    console.log('handleSiteIdUpdates', { v: value })
    onChange({ ...query, siteId: value });
    onRunQuery();
  }


  // this is bad do better
  useEffect(() => {
    if (!entity) {
      onChange({ ...query, entity: entity_options[0].value });
      onRunQuery();
    }

    updateSiteIds()
  }, [updateSiteIds])

  const { siteId, entity } = query;

  return (
    <>
      {/* <HorizontalGroup> */}
      <InlineField
        label="Category"
        grow
        labelWidth={20}
        tooltip="The category of the Netlify API to query for"
      >
        <Select
          options={entity_options}
          value={entity}
          defaultValue={entity_options[0]}
          allowCustomValue
          onChange={handleEntityChange}
        />
      </InlineField>
      {entities_requiring_site_id.includes(entity ?? '') && (
        <InlineField
          label="Site Id"
          grow
          labelWidth={20}
          tooltip="The Site Id to query By"
        >
          {/* <Select
            options={siteIdOptions}
            value={siteId}
            defaultValue={default_site_id}
            allowCustomValue
            onFocus={updateSiteIds}
            onChange={handleSiteIdUpdates}
          /> */}
          <Input name="importantInput" required value={siteId} onChange={(e) => {
            handleSiteIdUpdates(e.currentTarget.value)
          }} />
        </InlineField>

      )}
      {/* </HorizontalGroup> */}

      <HorizontalGroup>
        <ParametersEditor entity={entity} query={query} onChange={onChange} />
      </HorizontalGroup>

      <ParsingOptionsEditor query={query} onRunQuery={onRunQuery} data={data} onChange={onChange} editorType='query' actionConfig={{}} />
    </>
  );
}
