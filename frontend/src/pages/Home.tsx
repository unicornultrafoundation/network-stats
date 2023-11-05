import React, { useEffect, useState, useCallback } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import {
  PieChart, Pie, Cell, Tooltip,
  LabelList, Bar, BarChart, XAxis, YAxis
} from 'recharts';

import { scaleOrdinal } from 'd3-scale';
import { schemeCategory10 } from 'd3-scale-chromatic';

import { Grid, GridItem, useColorModeValue, Text, Center } from '@chakra-ui/react';
import { Card } from '../atoms/Card';
import { TooltipCard } from '../atoms/TooltipCard';
import {appendOtherGroup, getClients, getErrors} from '../data/DataMassager';
import { drilldownFilter, filterCount, FilterGroup, generateFilterGroupsFromQueryString, generateQueryStringFromFilterGroups } from '../data/FilterUtils';
import { Filtering } from '../organisms/Filtering';
import { Loader } from '../organisms/Loader';
import { knownNodesFilter, LayoutEightPadding, LayoutTwoColSpan, LayoutTwoColumn } from '../config';
import { CustomResponsiveContainer } from '../atoms/CustomResponsiveContainer';

const colors = scaleOrdinal(schemeCategory10).range();

interface NamedCount {
  name: string;
  count: number;
}

interface ClientData {
  versions: NamedCount[];
  versionsUnknown: number;
  clients: NamedCount[];
  errors: NamedCount[];
  clientsUnknown: number;
  operatingSystems: NamedCount[];
  operatingSystemsUnknown: number;
  languages: NamedCount[];
  languagesUnknown: number;
  countries: NamedCount[];
}

function Home() {
  const location = useLocation();
  const history = useHistory()
  const color = useColorModeValue("gray.800", "gray")
  const [data, setData] = useState<ClientData>()
  const [filters, setFilters] = useState<FilterGroup[]>()

  useEffect(() => {
    let searchFilters: FilterGroup[] = []
    if (location.search) {
      try {
        searchFilters = generateFilterGroupsFromQueryString(location.search);
      } catch (e) {
        console.error(e);
      }
    }

    if (filterCount(searchFilters) === 0) {
      // Deep clone since we are mutating it.
      searchFilters = JSON.parse(JSON.stringify(knownNodesFilter));
    }

    setFilters(searchFilters)
  }, [location.search])

  useEffect(() => {
    if (!filters) {
      return;
    }

    const run = async () => {
      const response = await fetch(`/v1/dashboard${generateQueryStringFromFilterGroups(filters)}`)
      const json: ClientData = await response.json()

      const [versions, unknownVersionsCount] = appendOtherGroup(json.versions)
      const [clients, unknownClientCount] = getClients(json.clients)
      const errors = getErrors(json.clients)
      const [languages, unknownLanguageCount] = appendOtherGroup(json.languages)
      const [operatingSystems, unknownOperatingSystemCount] = appendOtherGroup(json.operatingSystems)
      const [countries] = appendOtherGroup(json.countries)

      json.versions = versions
      json.versionsUnknown = unknownVersionsCount
      json.clients = clients
      json.errors = errors
      json.clientsUnknown = unknownClientCount
      json.languages = languages
      json.languagesUnknown = unknownLanguageCount
      json.operatingSystems = operatingSystems
      json.operatingSystemsUnknown = unknownOperatingSystemCount
      json.countries = countries;

      setData(json)
    }

    run()
  }, [filters])

  const onFiltersChanged = useCallback((filters: FilterGroup[]) => {
    history.push(location.pathname + generateQueryStringFromFilterGroups(filters))
  }, [history, location])

  const onClientClicked = useCallback((e: any) => e && onFiltersChanged(drilldownFilter(filters, 'name', e.activeLabel)), [filters, onFiltersChanged])
  const onOperatingSystemClicked = useCallback((e: any) => e && onFiltersChanged(drilldownFilter(filters, 'os_name', e.name)), [filters, onFiltersChanged])
  const onVersionClicked = useCallback((e: any) => e && onFiltersChanged(drilldownFilter(filters, 'version', e.activeLabel)), [filters, onFiltersChanged])

  if (!data) {
    return <Loader>Loading data...</Loader>
  }

  const renderLabelContent: React.FunctionComponent = (props: any) => {
    const { name, value, percent, x, y, midAngle } = props;
    return (
      <g transform={`translate(${x}, ${y})`} textAnchor={(midAngle < -90 || midAngle >= 90) ? 'end' : 'start'} fill={color}>
        <text x={0} y={0}>{`${name || "unknown"}`}</text>
        <text x={0} y={20}>{`${value} (${(percent * 100).toFixed(2)}%)`}</text>
      </g>
    );
  };

  const renderTooltipContent = (props: any): any => {
    if (!props.active || !props.payload || !props.payload.length) {
      return null
    }

    const { payload: { name, count } } = props.payload[0]
    return (
      <TooltipCard>
        <Text fontWeight="bold">{name}</Text>
        <Text>Count: {count}</Text>
      </TooltipCard>
    )
  };

  const barChartData = data.versions.length ? data.versions : data.clients
  const barChartTitle = data.versions.length ? 'Popular Versions' : 'Popular Clients'

  return (
    <Grid gridGap={LayoutEightPadding} templateColumns={LayoutTwoColumn} w="100%">
      <GridItem colSpan={LayoutTwoColSpan}>
        <Filtering filters={filters} onFiltersChange={onFiltersChanged} />
      </GridItem>


        <Card title={barChartTitle} contentHeight={barChartData.length * 40}>
          {barChartData.length === 0 && (
            <Center flex={1}>No data available</Center>
          )}
          {barChartData.length > 0 && (
            <CustomResponsiveContainer>
              <BarChart
                data={barChartData}
                layout="vertical"
                margin={{ left: 60, right: 30 }}
                onClick={data.versions.length ? onVersionClicked : onClientClicked}
              >
                <XAxis type="number" hide stroke={color} />
                <YAxis dataKey="name" type="category" interval={0} stroke={color} />
                <Tooltip cursor={false} content={renderTooltipContent} />
                <Bar dataKey="count">
                  {barChartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={colors[index % 10]} />
                  ))}
                  <LabelList position="right" />
                </Bar>
              </BarChart>
            </CustomResponsiveContainer>
          )}
        </Card>

      <Card title="Crawling Errors" contentHeight={data.errors.length * 40}>
        {data.errors.length === 0 && (
            <Center flex={1}>No data available</Center>
        )}
        {data.errors.length > 0 && (
            <CustomResponsiveContainer>
              <BarChart
                  data={data.errors}
                  layout="vertical"
                  margin={{ left: 400, right: 30 }}
                  onClick={data.versions.length ? onVersionClicked : onClientClicked}
              >
                <XAxis type="number" hide stroke={color} />
                <YAxis dataKey="name" type="category" interval={0} stroke={color} />
                <Tooltip cursor={false} content={renderTooltipContent} />
                <Bar dataKey="count">
                  {data.errors.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={colors[index % 10]} />
                  ))}
                  <LabelList position="right" />
                </Bar>
              </BarChart>
            </CustomResponsiveContainer>
        )}
      </Card>


      <Card title="Popular Operating Systems" contentHeight={300}>
        {data.operatingSystems.length === 0 && (
          <Center flex={1}>No data available</Center>
        )}
        {data.operatingSystems.length > 0 && (
          <CustomResponsiveContainer>
            <PieChart>
              <Pie
                data={data.operatingSystems}
                dataKey="count"
                startAngle={180}
                endAngle={-180}
                innerRadius={30}
                outerRadius={100}
                paddingAngle={data.operatingSystems.length === 1 ? 0 : 10}
                minAngle={data.operatingSystems.length === 1 ? 0 : 20}
                label={renderLabelContent}
                isAnimationActive={false}
                onClick={onOperatingSystemClicked}
                stroke="none"
              >
                {
                  data.operatingSystems.map((entry, index) => (
                    <Cell
                      key={`slice-${index}`}
                      fill={colors[index % 10] as string}
                    />
                  ))
                }
              </Pie>
            </PieChart>
          </CustomResponsiveContainer>
        )}

      </Card>

      <Card title="Popular Client Runtimes" contentHeight={300}>
        {data.languages.length === 0 && (
          <Center flex={1}>No data available</Center>
        )}
        {data.languages.length > 0 && (
          <CustomResponsiveContainer>
            <PieChart>
              <Pie
                data={data.languages}
                dataKey="count"
                startAngle={180}
                endAngle={-180}
                innerRadius={30}
                outerRadius={100}
                paddingAngle={data.languages.length === 1 ? 0 : 10}
                minAngle={data.languages.length === 1 ? 0 : 20}
                label={renderLabelContent}
                isAnimationActive={false}
                stroke="none"
              >
                {
                  data.languages.map((entry, index) => (
                    <Cell
                      key={`slice-${index}`}
                      fill={colors[index % 10] as string}
                    />
                  ))
                }
              </Pie>
            </PieChart>
          </CustomResponsiveContainer>
        )}
      </Card>
      <GridItem colSpan={LayoutTwoColSpan}>
        <Card title="Countries" contentHeight={data.countries.length * 40}>
          {data.countries.length === 0 && (
            <Center flex={1}>No data available</Center>
          )}
          {data.countries.length > 0 && (
            <CustomResponsiveContainer>
              <BarChart
                data={data.countries}
                layout="vertical"
                margin={{ left: 60, right: 30 }}
                onClick={data.versions.length ? onVersionClicked : onClientClicked}
              >
                <XAxis type="number" hide stroke={color} />
                <YAxis dataKey="name" type="category" interval={0} stroke={color} />
                <Tooltip cursor={false} content={renderTooltipContent} />
                <Bar dataKey="count">
                  {data.countries.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={colors[index % 10]} />
                  ))}
                  <LabelList position="right" />
                </Bar>
              </BarChart>
            </CustomResponsiveContainer>
          )}
        </Card>
      </GridItem>
    </Grid>
  );
}

export default Home;
