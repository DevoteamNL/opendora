package com.devoteam.backstagemock;

import com.jayway.jsonpath.Configuration;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;

import java.util.List;
import java.util.stream.Stream;

import static com.jayway.jsonpath.matchers.JsonPathMatchers.hasJsonPath;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class DeploymentFrequencyMetricTest {

    @Autowired
    private TestRestTemplate restTemplate;

    @ParameterizedTest
    @MethodSource("deploymentFrequencyMetricTestData")
    void mockDoraApiDeploymentFrequencyShouldReturnDataPoints(String type, String aggregation) {
        var dataPointMatcher = allOf(
                hasKey("key"),
                hasKey("value")
        );

        var url = "http://localhost:10666/dora/api/metric?type=%s&aggregation=%s".formatted(type, aggregation);

        var result = restTemplate.getForObject(url, String.class);

        var json = Configuration.defaultConfiguration().jsonProvider().parse(result);
        assertThat(json, hasJsonPath("$.aggregation", is(aggregation)));
        assertThat(json, hasJsonPath("$.dataPoints", isA(List.class)));
        assertThat(json, hasJsonPath("$.dataPoints[*]", hasItems(dataPointMatcher)));
        assertThat(json, hasJsonPath("$.dataPoints[*].key", hasItems(isA(String.class))));
        assertThat(json, hasJsonPath("$.dataPoints[*].value", hasItems(isA(Number.class))));
    }

    private Stream<Arguments> deploymentFrequencyMetricTestData() {
        return Stream.of(
                Arguments.of("df_average", "weekly"),
                Arguments.of("df_average", "monthly"),
                Arguments.of("df_average", "quarterly")
        );
    }

}
