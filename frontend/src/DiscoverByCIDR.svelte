<script lang="ts">
    import {
        Button,
        Column,
        ProgressBar,
        Row,
        StructuredList,
        StructuredListBody,
        StructuredListCell,
        StructuredListHead,
        StructuredListRow,
        StructuredListSkeleton,
        TextInput,
        Tile
    } from "carbon-components-svelte";
    import ArrowRight from "carbon-icons-svelte/lib/ArrowRight.svelte";
    import IPAddressRow from "./IPAddressRow.svelte";
    import {PingIPsByCIDR} from "../wailsjs/go/main/App";

    let loading: boolean = false;
    let inputCIDR: string | number = undefined;
    let discoveredIPs: { ip: string; id: string; port: number; lastPortStatus: boolean }[] = [];

    export let CIDR: string;


    function pingByCIDR(cidr: string | number): void {
        loading = true;
        cidr = cidr.toString(); // just to make TS stop complaining
        PingIPsByCIDR(cidr).then(result => {
            if (result === null) {
                loading = false;
                return [];
            }
            discoveredIPs = result.map((value: string) => {
                return {"ip": value, "id": value, "lastPortStatus": undefined, port: undefined}
            });
            loading = false;
        })
    }
</script>

<Row>
    <Column>
        <Row>
            <Tile>
                <TextInput type="text" value={CIDR} on:input={(event) => {inputCIDR = event.detail}}/>
            </Tile>
            <Tile>
                <Button kind="secondary" size="small" on:click={() => {pingByCIDR(inputCIDR)}}
                        icon={ArrowRight}
                        iconDescription="Ping">
                </Button>
            </Tile>
        </Row>
    </Column>
    <Column>
        <Tile>
            {#if loading}
                <ProgressBar helperText="Pinging..."/>
            {:else}
            {/if}
            <StructuredList condensed>
                <StructuredListHead>
                    <StructuredListRow head>
                        <StructuredListCell head>IP</StructuredListCell>
                        <StructuredListCell head>Port</StructuredListCell>
                        <StructuredListCell head>Ping</StructuredListCell>
                        <StructuredListCell head>Status</StructuredListCell>
                    </StructuredListRow>
                </StructuredListHead>
                <StructuredListBody>
                    {#each discoveredIPs as ipData}
                        <IPAddressRow ipData={ipData}/>
                    {/each}
                </StructuredListBody>
            </StructuredList>
            {#if loading}
                <StructuredListSkeleton rows={5}/>
            {:else}
            {/if}
        </Tile>
    </Column>
</Row>
