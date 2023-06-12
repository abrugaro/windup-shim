package conversion

import "strings"

// TODO: scope license file search to LICENSE files and use templates for tags to convert individual rules into one rule
const license_rules = `
- ruleID: discover-license
  description: "Discover project license"
  tag:
  - License={{matchingText}}
  when:
    or:
    - builtin.filecontent:
        pattern: "Apache License 1.0"
    - builtin.filecontent:
        pattern: "Apache License 1.1"
    - builtin.filecontent:
        pattern: "Apache License 2.0"
    - builtin.filecontent:
        pattern: "Mozilla Public License 2.0"
    - builtin.filecontent:
        pattern: "GNU GPL"
    - builtin.filecontent:
        pattern: "GNU LGPL"
    - builtin.filecontent:
        pattern: "CDDL"
    - builtin.filecontent:
        pattern: "Eclipse Public License 1.0"
    - builtin.filecontent:
        pattern: "BSD License"
    - builtin.filecontent:
        pattern: "Public Domain License"
`

// converted from https://github.com/windup/windup/blob/master/rules-java/api/src/main/java/org/jboss/windup/rules/apps/java/
// TODO: scope hard-coded IP to .java, pom.xml, .properties files
const java_rules = `
- ruleID: hardcoded-ip-address
  description: "Hardcoded IP Address\nWhen migrating environments, hard-coded IP addresses may need to be modified or eliminated."
  labels:
  - konveyor.io/target=cloud-readiness
  when:
    builtin.filecontent:
      pattern: ([0-9]{1,3}\.){3}[0-9]{1,3}
      filePattern: ".*\\.(java|properties)"
  category: mandatory
  effort: 1
  message: "When migrating environments, hard-coded IP addresses may need to be modified or eliminated."
- ruleID: discover-properties-file
  description: "Properties file"
  when:
    builtin.file:
      pattern: "^.*\\.properties$"
  tag: ["Properties"]
- ruleID: discover-manifest-file
  description: "Manifest file"
  when:
    builtin.file:
      pattern: "MANIFEST.MF"
  tag: ["Manifest"]
- ruleID: discover-java-files
  description: "Java source files"
  when:
    builtin.file:
      pattern: "*.java"
  tag: ["Java Source"]
- ruleID: discover-maven-xml
  description: "Maven XML file"
  when:
    builtin.file:
      pattern: "pom.xml"
  tag: ["Maven XML"]
`

const java_ee_rules = `
- ruleID: windup-discover-ejb-configuration
  tag: ["EJB XML"]
  when:
    builtin.xml:
      xpath: "/(jboss:ejb-jar or ejb-jar)"
- ruleID: windup-discover-spring-configuration
  tag: ["Spring XML"]
  labels:
  - "konveyor.io/fact=Spring Beans"
  when:
    builtin.xml:
      xpath: "/beans"
- ruleID: windup-discover-jpa-configuration
  labels:
  - "konveyor.io/fact=JPA"
  tag: ["JPA XML"]
  message: "Persistence unit"
  when:
    or:
      - builtin.xml:
          xpath: '/persistence[boolean(namespace-uri(/persistence)="http://java.sun.com/xml/ns/persistence")]'
      - builtin.xml:
          xpath: '/persistence[boolean(namespace-uri(/persistence)="http://xmlns.jcp.org/xml/ns/persistence")]'
      - builtin.xml:
          xpath: '/persistence[boolean(namespace-uri(/persistence)="https://jakarta.ee/xml/ns/persistence")]'
- ruleID: windup-discover-web-configuration
  tag: ["Web XML"]
  when:
    # TODO extract version as in rules-java-ee/addon/src/main/java/org/jboss/windup/rules/apps/javaee/rules/DiscoverWebXmlRuleProvider.java
    builtin.xml:
      xpath: /web-app
- ruleID: ejb-stateless-000
  labels:
  - "konveyor.io/fact=EJB"
  message: "Stateless session bean {{name}}"
  description: Stateless session beans
  when:
    or:
    - java.referenced:
        location: IMPLEMENTS_TYPE
        pattern: "(javax|jakarta).ejb.SessionBean"
    - java.referenced:
        location: ANNOTATION
        pattern: "(javax|jakarta).ejb.Stateless"
- ruleID: ejb-stateful-000
  labels:
  - "konveyor.io/fact=EJB"
  message: "Stateful session bean {{name}}"
  description: Stateful session beans
  when:
    java.referenced:
      location: ANNOTATION
      pattern: "(javax|jakarta).ejb.Stateful"
- ruleID: ejb-entity-000
  labels:
  - "konveyor.io/fact=EJB"
  message: "Entity bean {{name}}"
  description: "Entity beans"
  when:
    java.referenced:
      location: IMPLEMENTS_TYPE
      pattern: "(javax|jakarta).ejb.EntityBean"
- ruleID: ejb-message-driven-000
  labels:
  - "konveyor.io/fact=EJB"
  message: "Message driven bean {{name}}"
  description: "Message driven beans"
  when:
    or:
    - java.referenced:
        location: IMPLEMENTS_TYPE
        pattern: "(javax|jakarta).ejb.MessageDrivenBean"
    - java.referenced:
        location: ANNOTATION
        pattern: "(javax|jakarta).ejb.MessageDriven"
- ruleID: ejb-home-000
  labels:
  - "konveyor.io/fact=EJB"
  message: "EJB Home {{name}}"
  description: "EJB Homes"
  when:
    java.referenced:
      location: INHERITANCE
      pattern: "(javax|jakarta).ejb.(EJBHome|EJBLocalHome)"
- ruleID: ejb-object-000
  labels:
  - "konveyor.io/fact=EJB"
  message: "EJB Object {{name}}"
  description: "EJB Objects"
  when:
    java.referenced:
      location: INHERITANCE
      pattern: "(javax|jakarta).ejb.(EJBObject|EJBLocalObject)"
- ruleID: jpa-entity-000
  labels:
  - "konveyor.io/fact=JPA"
  message: "JPA Entity {{name}}"
  description: "JPA Entities"
  when:
    java.referenced:
      location: ANNOTATION
      pattern: "(javax|jakarta).persistence.Entity"
- ruleID: jpa-table-000
  labels:
  - "konveyor.io/fact=JPA"
  message: "JPA Table {{name}}"
  description: "JPA Tables"
  when:
    java.referenced:
      location: ANNOTATION
      pattern: "(javax|jakarta).persistence.Table"
- ruleID: jpa-query-000
  labels:
  - "konveyor.io/fact=JPA"
  message: "JPA Named Query {{name}}"
  description: "JPA Named Queries"
  when:
    java.referenced:
      location: ANNOTATION
      pattern: "(javax|jakarta).persistence.(NamedQuery|NamedQueries)"
`

func GetDiscoveryRules() string {
	rules := []string{}
	rules = append(rules, strings.Trim(license_rules, "\n"))
	rules = append(rules, strings.Trim(java_rules, "\n"))
	rules = append(rules, strings.Trim(java_ee_rules, "\n"))
	return strings.Join(rules, "\n")
}

func GetDiscoveryRuleset() string {
	return `
name: discovery-rules
labels:
- discovery
`
}
