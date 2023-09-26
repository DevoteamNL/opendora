import {
  Content,
  ContentHeader,
  Header,
  Page,
  Progress,
  ResponseErrorPanel,
  SupportButton,
} from '@backstage/core-components';
import { getEntityRelations, useEntity } from '@backstage/plugin-catalog-react';
import { Grid } from '@material-ui/core';
import React, { useEffect } from 'react';
import { DeploymentFrequencyData } from '../../models/DeploymentFrequencyData';
import GroupDataService from '../../services/GroupDataService';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';
import { DropdownComponent } from '../DropdownComponent/DropdownComponent';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';
import './DashboardComponent.css';
import { calculateAverage } from '../../utils/helpers';

export const DashboardComponent = () => {
  const [chartData, setChartData] =
    React.useState<DeploymentFrequencyData | null>(null);

  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('weekly');
  const [deploymentAverage, setDeploymentAverage] = React.useState(0);

  const { entity } = useEntity();
  const groupName = getEntityRelations(entity, 'ownedBy')[0]?.name;

  const [dataError, setDataError] = React.useState<Error | undefined>(
    undefined,
  );

  useEffect(() => {
    if (!groupName) return;
    GroupDataService.getMockData(groupName, selectedTimeUnit).then(
      response => {
        setChartData(response);
        setDeploymentAverage(calculateAverage(response));
      },
      (error: Error) => {
        setDataError(error);
      },
    );
  }, [groupName, selectedTimeUnit]);
  const chartOrProgressComponent = chartData ? (
    <BarChartComponent deploymentFrequencyData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

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
                  title={selectedTimeUnit.toUpperCase() + " average number of deployments"}
                  highlight={deploymentAverage.toString()}
                  textColour="positiveHighlight"
                />
              </div>
            </Grid>
            <Grid item xs={6} className="gridBorder">
              <div className="gridBoxText">
                <HighlightTextBoxComponent
                  title="Placeholder"
                  highlight="Placeholder"
                  text="*Placeholder info"
                  textColour="warning"
                />
              </div>
            </Grid>

            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                {dataError ? (
                  <ResponseErrorPanel error={dataError} />
                ) : (
                  chartOrProgressComponent
                )}
              </div>
            </Grid>
          </Grid>

          <Grid item />
        </Grid>
      </Content>
    </Page>
  );
};
