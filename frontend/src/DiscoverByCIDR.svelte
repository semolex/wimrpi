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
    import {AutoDetect, PingIPsByCIDR} from "../wailsjs/go/main/App";
    import {onMount} from "svelte";

    let loading: boolean = false;
    let inputCIDR: string | number = undefined;
    let discoveredIPs: { ip: string; id: string; port: number; lastPortStatus: boolean }[] = [];

    function isValidCIDR(cidr: string): boolean {
        const pattern = new RegExp(
            '^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])' +
            '\\/(\\d|[1-2]\\d|3[0-2])$'
        );

        return pattern.test(cidr);
    }

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

    function getOwnCIDR(): void {
        AutoDetect().then(result => inputCIDR = result)
    }

    onMount(async () => {
        getOwnCIDR();
    });
</script>

<Row>
    <Column>
        <Row style="margin-left: 0rem;">
            <Column>
                <Row>
                    <TextInput
                            type="text" bind:value={inputCIDR}
                            invalid={inputCIDR && !isValidCIDR(inputCIDR.toString())}
                            invalidText="Invalid CIDR"
                    />
                    <Button kind="secondary" size="small" on:click={() => {pingByCIDR(inputCIDR)}}
                            icon={ArrowRight}
                            disabled={!inputCIDR || !isValidCIDR(inputCIDR.toString())}
                            iconDescription="Ping">
                    </Button>
                </Row>
            </Column>
            <Column>
                <Row>
                </Row>
            </Column>
        </Row>
        <Tile>
            {#if loading}
                <ProgressBar helperText="Pinging..."/>
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
            {/if}
        </Tile>
    </Column>
</Row>
