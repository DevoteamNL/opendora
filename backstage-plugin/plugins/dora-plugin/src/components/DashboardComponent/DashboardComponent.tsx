import React, { useEffect } from 'react';
import { Grid } from '@material-ui/core';
import {
  Header,
  Page,
  Content,
  ContentHeader,
  SupportButton,
  Progress,
  ResponseErrorPanel,
} from '@backstage/core-components';
import './DashboardComponent.css';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';
import { DropdownComponent } from '../DropdownComponent/DropdownComponent';

import GroupDataService from '../../services/GroupDataService';
import { DeploymentFrequencyData } from '../../models/DeploymentFrequencyData';
import { useEntity } from '@backstage/plugin-catalog-react';
import type { EntityRelation } from '@backstage/catalog-model';

export const DashboardComponent = () => {
  const [chartData, setChartData] =
    React.useState<DeploymentFrequencyData | null>(null);

  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('Weekly');

  const entity = useEntity();
  const groupOwnerRelationship = entity.entity.relations?.find(
    relation => relation.type === 'ownedBy',
  );

  const [dataError, setDataError] = React.useState<Error | undefined>(
    undefined,
  );

  useEffect(() => {
    if (!groupOwnerRelationship) return;
    // TODO: Create a fix for this in backstage repo
    const group = (
      groupOwnerRelationship as EntityRelation & {
        target: { name: string; kind: string; namespace: string };
      }
    ).target.name;
    GroupDataService.getMockData(group, selectedTimeUnit).then(
      response => {
        setChartData(response);
      },
      (error: Error) => {
        setDataError(error);
      },
    );
  }, [groupOwnerRelationship, selectedTimeUnit]);

  return (
    <Page themeId="tool">
      <Header
        title="Devoteam DORA plugin"
        subtitle="Through insight to perfection"
      />
      <Content>
        <ContentHeader title="DORA metrics">
          <SupportButton>Plugin for displaying DORA metrics</SupportButton>
        </ContentHeader>
        <Grid container spacing={3} direction="column">
          <Grid container>
            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                <h1>Deployment statistics</h1>
                <p>Analysis of successful deployments and CFR</p>
                <Grid container>
                  <Grid item xs={4}>
                    <DropdownComponent
                      onSelect={setSelectedTimeUnit}
                      selection={selectedTimeUnit}
                    />
                  </Grid>
                </Grid>
              </div>
            </Grid>

            <Grid item xs={6} className="gridBorder">
              <div className="gridBoxText">
                <HighlightTextBoxComponent
                  title="Average number of deployments per week"
                  highlight="31"
                  textColour="positiveHighlight"
                />
              </div>
            </Grid>
            <Grid item xs={6} className="gridBorder">
              <div className="gridBoxText">
                <HighlightTextBoxComponent
                  title="Overall change failure rate"
                  highlight="5.2%"
                  text="*calculated on failures and incidents"
                  textColour="warning"
                />
              </div>
            </Grid>

            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                {chartData ? (
                  <BarChartComponent deploymentFrequencyData={chartData} />
                ) : (
                  <Progress variant="indeterminate" />
                )}

                {dataError ? <ResponseErrorPanel error={dataError} /> : null}
              </div>
            </Grid>
          </Grid>

          <Grid item />
        </Grid>
      </Content>
    </Page>
  );
};
