import { InlineField, Input } from "@grafana/ui"
import { QueryOptionGroup } from "./QueryOptionsGroup"
import React from "react"

type Props = {

}

export const ParametersEditor = (props: Props) => {

    return (
        <QueryOptionGroup title="Optional Parameters">
            <InlineField label="Limit" labelWidth={20}
                tooltip="Name of the field in the response to read data from"
                grow>
                <Input name="hello" />
            </InlineField>
        </QueryOptionGroup>
    )
}