<script lang="ts">
    import 'carbon-components-svelte/css/white.css'
    import {Content, Grid, Header, SkipToContent, Tab, TabContent, Tabs,} from "carbon-components-svelte";
    import DiscoverByCIDR from "./DiscoverByCIDR.svelte";
    import FindOpenPorts from "./FindOpenPorts.svelte";
</script>

<Header company="WIMRPI" platformName="Ping it">
    <svelte:fragment slot="skip-to-content">
        <SkipToContent/>
    </svelte:fragment>
</Header>

<Content>
    <Grid>
        <Tabs autoWidth>
            <Tab label="Discover by CIDR"/>
            <Tab label="Find open ports"/>
            <svelte:fragment slot="content">
                <TabContent>
                    <DiscoverByCIDR/>
                </TabContent>
                <TabContent>
                    <FindOpenPorts/>
                </TabContent>
            </svelte:fragment>
        </Tabs>
    </Grid>
</Content>





