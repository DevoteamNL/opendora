package com.devoteam.backstagemock;

import com.jayway.jsonpath.Configuration;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.test.context.ActiveProfiles;

import java.util.List;
import java.util.stream.Stream;

import static com.jayway.jsonpath.matchers.JsonPathMatchers.hasJsonPath;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;

@ActiveProfiles("test")
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class DeploymentFrequencyMetricTest {

    @Value("${wiremock.server.port}")
    private int wireMockPort;

    @Autowired
    private TestRestTemplate restTemplate;

    @ParameterizedTest
    @MethodSource("deploymentFrequencyMetricTestData")
    void mockDoraApiDeploymentFrequencyShouldReturnDataPoints(String type, String aggregation) {
        var dataPointMatcher = allOf(
                hasKey("key"),
                hasKey("value")
        );

        var url = "http://localhost:%s/dora/api/metric?type=%s&aggregation=%s".formatted(wireMockPort, type, aggregation);

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
                Arguments.of("df_average", "quarterly"),
                Arguments.of("df_count", "weekly"),
                Arguments.of("df_count", "monthly"),
                Arguments.of("df_count", "quarterly")
        );
    }

}
