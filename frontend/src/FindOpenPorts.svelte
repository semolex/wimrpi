<script type="ts">
    import {GeneratePortTreeForIP} from "../wailsjs/go/main/App";
    import type {main} from "../wailsjs/go/models";
    import TrashCan from "carbon-icons-svelte/lib/TrashCan.svelte";
    import {Button, Column, Row, StructuredListSkeleton, TextInput, Tile, TreeView} from "carbon-components-svelte";
    import ArrowRight from "carbon-icons-svelte/lib/ArrowRight.svelte";
    import Add from "carbon-icons-svelte/lib/Add.svelte";
    import type {Writable} from "svelte/store";
    import {writable} from "svelte/store";

    let ips: Writable<string[]> = writable([""]);
    let portTree: main.PortTree[] = [];
    let loading = false;

    function generatePortTree(ip: string | number) {
        ip = ip.toString();
        loading = true;
        GeneratePortTreeForIP(ip).then((result) => {
            if (result !== null) {
                portTree = portTree.concat(result);
            }
            loading = false;
        })
            .catch(error => {
                console.error(error);
                loading = false;
            });
    }

    function addIp() {
        $ips = [...$ips, ""];
    }

    function removeIp(index: number) {
        let updatedIps = $ips.filter((_, i) => i !== index);
        ips.set([]); // Force an update
        ips.set(updatedIps); // Assign the updated IPs
    }

    function isValidIP(ip: string): boolean {
        const pattern = new RegExp(
            '^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.' +
            '(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.' +
            '(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.' +
            '(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$'
        );

        return pattern.test(ip);
    }


</script>
<Row>
    <Column>
        {#each $ips as ip, index (index)}
            <Row>
                <Tile>
                    <TextInput
                            type="text"
                            bind:value={ip}
                            invalid={ip !== "" && !isValidIP(ip)}
                            invalidText="Invalid IP."
                    />
                </Tile>
                <Tile>
                    <Button kind="secondary" size="small" on:click={addIp}
                            icon={Add}
                            iconDescription="Add">
                    </Button>
                    <Button kind="danger" size="small" on:click={() => {
                        removeIp(index);
                    }}
                            icon={TrashCan}
                            disabled={$ips.length === 1 || index === 0}
                            iconDescription="Remove">
                    </Button>
                    <Button kind="secondary" size="small" on:click={() => {generatePortTree(ip)}}
                            icon={ArrowRight}
                            disabled={ip === "" || !isValidIP(ip)}
                            iconDescription="Ping">
                    </Button>
                </Tile>
            </Row>
        {/each}

    </Column>
    <Column>
        <Tile>
            {#if loading}
                <StructuredListSkeleton/>
            {:else}
                <TreeView children={portTree} labelText="Open ports"/>
            {/if}
        </Tile>
    </Column>
</Row>