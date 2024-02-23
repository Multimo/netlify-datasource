import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput, useStyles2 } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, GrafanaTheme2 } from '@grafana/data';
import { NetlifyDataSourceOptions, NetlifySecureJsonData } from '../types';
import { DataSourceDescription, ConfigSection } from '@grafana/experimental';
import { css } from '@emotion/css';

interface Props extends DataSourcePluginOptionsEditorProps<NetlifyDataSourceOptions> { }

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const styles = useStyles2(getStyles);

  const onSiteIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      siteId: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  const onAccountIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      accountId: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };


  // Secure field (only sent to the backend)
  const onAPIKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        accessToken: event.target.value,
      },
    });
  };

  const onResetAPIKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        accessToken: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        accessToken: '',
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as NetlifySecureJsonData;

  return (
    <div className="gf-form-group">
      <DataSourceDescription dataSourceName='Netlify Datasource' docsLink='' />

      <hr className={styles.break} />

      <ConfigSection
        title="Authentication"
        description="Provide the Access Token the grafana plugin will use to authenicate with the Netlify Api. "
      >
        <InlineField label="Access Token" labelWidth={20} tooltip="Netlify Access token found in the User Settings -> Applications in Netlify">
          <SecretInput
            isConfigured={(secureJsonFields && secureJsonFields.accessToken) as boolean}
            value={secureJsonData.accessToken || ''}
            placeholder="Access Token"
            width={40}
            required
            onReset={onResetAPIKey}
            onChange={onAPIKeyChange}
          />
        </InlineField>
      </ConfigSection>

      <hr className={styles.break} />

      <ConfigSection
        title="Additional settings"
        description="Defaults the plugin will use in the query params"
        isCollapsible
        isInitiallyOpen={false}
      >
        <InlineField label="Default Site ID" labelWidth={20}>
          <Input
            onChange={onSiteIdChange}
            value={jsonData.siteId || ''}
            placeholder="Your Site ID"
            width={40}
          />
        </InlineField>
        <InlineField label="Account ID" labelWidth={20}>
          <Input
            onChange={onAccountIdChange}
            value={jsonData.accountId || ''}
            placeholder="Your Acount ID"
            width={40}
          />
        </InlineField>
      </ConfigSection>
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    break: css({
      marginTop: theme.spacing(4),
      marginBottom: theme.spacing(4)
    })
  }
}