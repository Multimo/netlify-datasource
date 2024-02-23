import { css } from '@emotion/css';
import React, { useState } from 'react';
import { GrafanaTheme2 } from '@grafana/data';
import { Collapse, HorizontalGroup, useStyles2 } from '@grafana/ui';

type Props = {
    title: string;
    description?: string;
    defaultIsOpen?: boolean;
    children: React.ReactNode;
};
export function QueryOptionGroup({ title, description, defaultIsOpen = false, children }: Props) {
    const [isOpen, setIsOpen] = useState(defaultIsOpen);
    const styles = useStyles2(getStyles);

    return (
        <div className={styles.wrapper}>
            <Collapse
                className={styles.collapse}
                collapsible
                isOpen={isOpen}
                onToggle={() => setIsOpen(!isOpen)}
                label={
                    <HorizontalGroup>
                        <h6 className={styles.title}>{title}</h6>
                        {!isOpen && description && (
                            <div className={styles.description}>
                                <span>{description}</span>
                            </div>
                        )}
                    </HorizontalGroup>
                }
            >
                <div className={styles.body}>{children}</div>
            </Collapse>
        </div>
    );
}

const getStyles = (theme: GrafanaTheme2) => {
    return {
        collapse: css({
            backgroundColor: 'unset',
            border: 'unset',
            marginBottom: 0,

            ['> button']: {
                padding: theme.spacing(0, 1),
            },
        }),
        wrapper: css({
            width: '100%',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'baseline',
            margin: theme.spacing(1, 0),
        }),
        title: css({
            flexGrow: 1,
            overflow: 'hidden',
            fontSize: theme.typography.bodySmall.fontSize,
            fontWeight: theme.typography.fontWeightMedium,
            margin: 0,
        }),
        description: css({
            color: theme.colors.text.secondary,
            fontSize: theme.typography.bodySmall.fontSize,
            fontWeight: theme.typography.bodySmall.fontWeight,
            paddingLeft: theme.spacing(2),
            gap: theme.spacing(2),
            display: 'flex',
        }),
        body: css({
            paddingTop: theme.spacing(0.5),
        }),
        stats: css({
            margin: '0px',
            color: theme.colors.text.secondary,
            fontSize: theme.typography.bodySmall.fontSize,
        }),
        tooltip: css({
            marginRight: theme.spacing(0.25),
        }),
    };
};
