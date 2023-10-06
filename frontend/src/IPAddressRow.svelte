<script lang="ts">
    import {
        Button, TextInput, StructuredListCell, StructuredListRow, Tag
    } from "carbon-components-svelte";
    import ConnectTarget from "carbon-icons-svelte/lib/ConnectTarget.svelte";
    import CheckmarkFilled from "carbon-icons-svelte/lib/CheckmarkFilled.svelte";
    import ErrorFilled from "carbon-icons-svelte/lib/ErrorFilled.svelte";
    import Unknown from "carbon-icons-svelte/lib/Unknown.svelte";
    import {PingPort} from "../wailsjs/go/main/App";

    export let ipData: { ip: string; id: string; port: number; lastPortStatus: boolean };

    function pingPort() {
        PingPort(ipData.ip, ipData.port).then(result => {
            ipData.lastPortStatus = result;
        });
    }
</script>

<StructuredListRow>
    <StructuredListCell>{ipData.ip}</StructuredListCell>
    <StructuredListCell style="width: 20%">
        <TextInput style="width: 45%" type="number" size="sm" hideLabel on:input={
            (event) => { ipData.port = +event.detail; }}/>
    </StructuredListCell>
    <StructuredListCell>
        <Button size="small" kind="tertiary" iconDescription="Ping"
                disabled={ipData.port === undefined || ipData.port === null}
                icon={ConnectTarget}
                on:click={() => {
                    pingPort();
                }}/>
    </StructuredListCell>
    <StructuredListCell>
        {#if ipData.lastPortStatus === undefined}
            <Tag icon={Unknown} type="gray"></Tag>
        {:else if ipData.lastPortStatus === true}
            <Tag icon={CheckmarkFilled} type="green"></Tag>
        {:else if ipData.lastPortStatus === false}
            <Tag icon={ErrorFilled} type="red"></Tag>
        {:else}
            <Tag icon={Unknown} type="gray"></Tag>
        {/if}
    </StructuredListCell>
</StructuredListRow>
