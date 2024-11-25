# 🖖 CloudStatus.Page

> [!NOTE]
> Archived version of [cloudstatus.page](https://cloudstatus.page) can be found at https://dn0.github.io/CloudStatusPage/

## Why

The goal of this project is simple: to bring greater transparency to the operational status of major public cloud services (and potentially other SaaS platforms). In today’s software-driven world, so much depends on the reliability of big cloud providers. While each of them offers a status page, it often paints a picture that looks perfect — perhaps a little too perfect.

## How it works

1. **Monitoring Agents:** A simple monitoring agent is deployed in every region of each cloud (AWS, Azure, GCP).
2. **Periodic Probes:** These agents regularly run "monitoring probes" - small tests that interact with various managed services in each cloud region.
3. **Official SDKs:** Monitoring probes are simply calling cloud's APIs using their official SDKs.
4. **Regional Testing:** All tests are performed within the same region as the monitoring agent so that reported latencies are purely regional.
5. **Data Collection:** Results from these probes are asynchronously collected in a central database and displayed on this site.
6. **Continuous Analysis:** Collected data are constantly analyzed for potential latency issues and probe failures. If something seems off, the system generates alerts, which could escalate into incidents (once verified).
