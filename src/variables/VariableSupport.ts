import { CustomVariableSupport, DataQueryRequest, DataQueryResponse, VariableSupportType } from '@grafana/data';

import { DataSource } from 'datasource';
import { QueryEditor } from '../components/QueryEditor';
import { NetlifyQuery, NetlifyDataSourceOptions } from 'types';
import { Observable, from, map } from 'rxjs';



export class VariableSupport extends CustomVariableSupport<DataSource, NetlifyQuery, NetlifyQuery, NetlifyDataSourceOptions> {
    constructor(private readonly datasource: DataSource) {
        super();
        this.datasource = datasource;

        this.query = this.query.bind(this);
    }

    query(request: DataQueryRequest<NetlifyQuery>): Observable<DataQueryResponse> {
        console.log("TestDataVariableSupport.query", request)
        // const executeObservable = from(this.ds.query(request.targets[0]));
        return from(this.datasource.getSiteIds()).pipe(
            map((data) => {
                const results = data.map((site) => ({ text: site, value: site }))
                console.log({ data })
                return { data: results }
            })
        )
    }

    getType(): VariableSupportType {
        return VariableSupportType.Custom;
    }

    editor = QueryEditor
}
