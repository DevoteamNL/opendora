package com.devoteam.backstagemock.config;

import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.client.WireMock;
import lombok.SneakyThrows;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;

import java.nio.file.Files;

@Configuration
class WireMockConfiguration {

    @Bean
    WireMockServer wireMockServer(@Value("${wiremock.server.port}") int wireMockPort, WireMockStubConfigurationProperties properties) {
        var wireMockServer = new WireMockServer(wireMockPort);

        properties.stubs().forEach(stub -> wireMockServer.stubFor(
                        WireMock.get(WireMock.urlPathEqualTo(stub.path()))
                                .withQueryParams(stub.queryParameters())
                                .willReturn(WireMock.aResponse()
                                        .withStatus(HttpStatus.OK.value())
                                        .withHeader(HttpHeaders.CONTENT_TYPE, stub.mediaType())
                                        .withBody(body(stub.stubFile())))
                )
        );

        wireMockServer.start();

        return wireMockServer;
    }

    @SneakyThrows
    private String body(String fileName) {
        var file = new ClassPathResource(fileName).getFile();

        return Files.readString(file.toPath());
    }

}
