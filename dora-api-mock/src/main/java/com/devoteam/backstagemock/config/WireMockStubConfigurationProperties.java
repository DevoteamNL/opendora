package com.devoteam.backstagemock.config;

import com.github.tomakehurst.wiremock.client.WireMock;
import com.github.tomakehurst.wiremock.matching.StringValuePattern;
import org.springframework.boot.context.properties.ConfigurationProperties;

import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

@ConfigurationProperties(prefix = "dora-api-mock")
record WireMockStubConfigurationProperties(List<Stubs> stubs) {

    record Stubs(String name, String path, String stubFile, String mediaType, Map<String, String> parameters) {

        public Map<String, StringValuePattern> queryParameters() {
            var params = Optional.ofNullable(parameters).orElse(Map.of());

            return params.entrySet()
                    .stream()
                    .collect(Collectors.toMap(Map.Entry::getKey, e -> WireMock.matching(e.getValue())));
        }

    }

}
