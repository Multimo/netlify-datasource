import React, { useEffect } from 'react';
import { css } from '@emotion/css';
import { InlineField, CodeEditor, Select, Icon, useTheme2 } from '@grafana/ui';
import { SelectableValue, toOption } from '@grafana/data';
import { QueryOptionGroup } from './QueryOptionsGroup';
// import { ActionConfig } from '../../types';
// import { EditorType } from './types';

type ActionConfig = any

import { NetlifyQuery } from 'types';
type Props = {
    query: NetlifyQuery;
    actionConfig?: ActionConfig;
    onChange: (query: NetlifyQuery) => void;
    editorType: 'query' | 'variable';
    data: any,
};

export const ParsingOptionsEditor = ({ query, actionConfig, onRunQuery, onChange, editorType, data }: Props) => {
    console.log('ParsingOptionsEditor', { query, actionConfig, onChange, editorType, data })
    const refId = query.refId;
    const series = data.series.find((s: any) => s.refId === refId);

    const onColumnsFilterChange = (options: SelectableValue<string>[]) => {
        const values = options.map((o) => o.value!);
        onChange({
            ...query,
            parsingOptions: {
                selectedFields: values
            },
        });
        onRunQuery()
    };
    const theme = useTheme2();
    const styles = {
        variableInfo: css({
            display: 'flex',
            alignItems: 'top',
            ...theme.typography.bodySmall,
            color: theme.colors.text.secondary,
        }),
        variableInfoIcon: css({
            margin: theme.spacing(1 / 8, 0.5, 0, 0),
        }),
    };

    const entity = query.entity;
    useEffect(() => {
        // reset field if entity changes
        onChange({
            ...query,
            parsingOptions: {
                selectedFields: []
            },
        });
    }, [entity])

    return (
        <>
            <QueryOptionGroup title="Transformation options" defaultIsOpen={!!query.parsingOptions?.selectedFields.length}>
                <InlineField
                    label="Select fields"
                    labelWidth={20}
                    tooltip="When filled in only selected fields will be returned"
                    grow
                >
                    <Select
                        options={series?.fields.map((f: any) => toOption(f.name))}
                        value={query.parsingOptions?.selectedFields}
                        isMulti
                        allowCustomValue
                        onChange={(options) => onColumnsFilterChange(options as SelectableValue<string>[])}
                    />
                </InlineField>
            </QueryOptionGroup>
            {editorType === 'variable' && (
                <p className={styles.variableInfo}>
                    <Icon name="info-circle" className={styles.variableInfoIcon} />
                    <span>
                        Note: first two selected fields will be used as variable label and value respectively. If only one field is
                        selected it will be used for both variable label and value.
                    </span>
                </p>
            )}
        </>
    );
};
