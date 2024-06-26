plugins {
    kotlin("multiplatform") version "2.0.0"
}

group = "me.kolterdyx"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

kotlin {
    linuxX64("linux") {
        binaries {
            executable("mcbasic")
        }
    }
    mingwX64("windows") {
        binaries {
            executable("mcbasic")
        }
    }
    applyDefaultHierarchyTemplate()
    sourceSets {
        commonMain.dependencies {
            implementation("com.github.ajalt.clikt:clikt:4.4.0")
        }
    }
}

tasks.withType<Wrapper> {
    gradleVersion = "8.5"
    distributionType = Wrapper.DistributionType.BIN
}