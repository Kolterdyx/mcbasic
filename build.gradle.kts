plugins {
    kotlin("multiplatform") version "2.0.0"
    id("io.kotest.multiplatform") version "5.9.1"
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
        commonTest.dependencies {
            implementation("io.kotest:kotest-framework-engine:5.9.1")
            implementation("io.kotest:kotest-assertions-core:5.9.1")
            implementation("io.kotest:kotest-property:5.9.1")
        }
    }
}

tasks.withType<Wrapper> {
    gradleVersion = "8.5"
    distributionType = Wrapper.DistributionType.BIN
}