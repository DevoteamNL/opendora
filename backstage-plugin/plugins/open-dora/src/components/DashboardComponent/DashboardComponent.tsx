import {
  Content,
  Header,
  Page,
  Progress,
  ResponseErrorPanel,
  SupportButton,
} from '@backstage/core-components';
import { useApi } from '@backstage/core-plugin-api';
import { getEntityRelations, useEntity } from '@backstage/plugin-catalog-react';
import { Grid } from '@material-ui/core';
import React, { useEffect, useReducer } from 'react';
import { MetricData } from '../../models/MetricData';
import { groupDataServiceApiRef } from '../../services/GroupDataService';
import { BarChartComponent } from '../BarChartComponent/BarChartComponent';
import { DropdownComponent } from '../DropdownComponent/DropdownComponent';
import './DashboardComponent.css';
import { ChartErrors } from '../../models/CustomErrors';
import { dfBenchmarkKey } from '../../models/DfBenchmarkData';
import { HighlightTextBoxComponent } from '../HighlightTextBoxComponent/HighlightTextBoxComponent';
import '../../i18n';
import { useTranslation } from 'react-i18next';

export interface DashboardComponentProps {
  entityName?: string;
  entityGroup?: string;
}

function dataErrorReducer(
  currentErrors: ChartErrors,
  action: { type: keyof ChartErrors; error: Error },
) {
  return {
    ...currentErrors,
    [action.type]: action.error,
  };
}

export const DashboardComponent = ({
  entityName,
  entityGroup,
}: DashboardComponentProps) => {
  // Overview
  const [dfOverview, setDfOverview] = React.useState<dfBenchmarkKey | null>(
    null,
  );

  const [t] = useTranslation();

  // Charts
  const [chartData, setChartData] = React.useState<MetricData | null>(null);
  const [chartDataAverage, setChartDataAverage] =
    React.useState<MetricData | null>(null);
  const [selectedTimeUnit, setSelectedTimeUnit] = React.useState('weekly');

  const initialErrors: ChartErrors = {
    countError: null,
    averageError: null,
    dfBenchmarkError: null,
  };

  const [dataError, dispatch] = useReducer(dataErrorReducer, initialErrors);

  const groupDataService = useApi(groupDataServiceApiRef);

  useEffect(() => {
    groupDataService
      .retrieveMetricDataPoints({
        type: 'df_count',
        team: entityGroup,
        aggregation: selectedTimeUnit,
        project: entityName,
      })
      .then(
        response => {
          setChartData(response);
        },
        (error: Error) => {
          dispatch({
            type: 'countError',
            error: error,
          });
        },
      );

    groupDataService
      .retrieveMetricDataPoints({
        type: 'df_average',
        team: entityGroup,
        aggregation: selectedTimeUnit,
        project: entityName,
      })
      .then(
        response => {
          setChartDataAverage(response);
        },
        (error: Error) => {
          dispatch({
            type: 'averageError',
            error: error,
          });
        },
      );

    groupDataService.retrieveBenchmarkData({ type: 'df' }).then(
      response => {
        setDfOverview(response.key);
      },
      (error: Error) => {
        dispatch({
          type: 'dfBenchmarkError',
          error: error,
        });
      },
    );
  }, [entityGroup, entityName, selectedTimeUnit, groupDataService]);

  const chartOrProgressComponent = chartData ? (
    <BarChartComponent metricData={chartData} />
  ) : (
    <Progress variant="indeterminate" />
  );

  const chartOrProgressComponentAverage = chartDataAverage ? (
    <BarChartComponent metricData={chartDataAverage} />
  ) : (
    <Progress variant="indeterminate" />
  );

  return (
    <Page themeId="tool">
      <Header
        title="OpenDORA (by Devoteam)"
        subtitle="Through insight to perfection"
      >
        <SupportButton>Plugin for displaying DORA Metrics</SupportButton>
      </Header>
      <Content>
        <Grid container spacing={3} direction="column">
          <Grid container>
            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
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

            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                <Grid container>
                  <Grid item xs={3}>
                    <HighlightTextBoxComponent
                      title=""
                      text=""
                      highlight={
                        dfOverview
                          ? t(
                              `deployment_frequency.overall_labels.${dfOverview}`,
                            )
                          : t('custom_errors.data_unavailable')
                      }
                      // to do: think of text colouring for different scenarios
                      textColour="positiveHighlight"
                    />
                  </Grid>
                </Grid>
              </div>
            </Grid>

            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                <h1>{t('deployment_frequency.labels.deployment_frequency')}</h1>
                {dataError.countError ? (
                  <ResponseErrorPanel error={dataError.countError} />
                ) : (
                  chartOrProgressComponent
                )}
              </div>
            </Grid>
            <Grid item xs={12} className="gridBorder">
              <div className="gridBoxText">
                <h1>
                  {t(
                    'deployment_frequency.labels.deployment_frequency_average',
                  )}
                </h1>
                {dataError.averageError ? (
                  <ResponseErrorPanel error={dataError.averageError} />
                ) : (
                  chartOrProgressComponentAverage
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

export const EntityDashboardComponent = () => {
  const { entity } = useEntity();
  const entityGroup = getEntityRelations(entity, 'ownedBy')[0]?.name;
  const entityName = entity.metadata.name;

  return (
    <DashboardComponent entityName={entityName} entityGroup={entityGroup} />
  );
};
