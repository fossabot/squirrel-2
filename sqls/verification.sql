/* Verify nep5 addresses */
select
    t.asset_id,
    t.symbol,
    t.addresses,
    t.count_from_addr_asset
from
    (select
         nep5.asset_id,
         nep5.symbol,
         nep5.addresses,
         count(distinct addr_asset.address) count_from_addr_asset
     from nep5
         join addr_asset on nep5.asset_id = addr_asset.asset_id
     group by nep5.asset_id) t
where t.count_from_addr_asset != t.addresses;

/* Verify nep5 holding_addresses */
select
    t.asset_id,
    t.symbol,
    t.holding_addresses,
    t.count_from_addr_asset
from (
         select
             nep5.asset_id,
             nep5.symbol,
             nep5.holding_addresses,
             count(distinct addr_asset.address) count_from_addr_asset
         from nep5
             join addr_asset on nep5.asset_id = addr_asset.asset_id
         where addr_asset.balance > 0
         group by nep5.asset_id) t
where t.count_from_addr_asset != t.holding_addresses;

/* Verify asset addresses */
select
    t.asset_id,
    t.name,
    t.addresses,
    t.count_from_addr_asset
from
    (
        select
            asset.asset_id,
            asset.name,
            asset.addresses,
            count(distinct addr_asset.address) count_from_addr_asset
        from addr_asset
            join asset on addr_asset.asset_id = asset.asset_id
        group by asset.asset_id) t
where t.addresses != t.count_from_addr_asset;

/* Verify total addresses */
select count(address) as addr_from_address
from address
union all
select count(distinct address) as _addr_from_addr_asset
from addr_asset;

/* Verify total supply of nep5 */
SELECT
    t.asset_id,
    t.symbol,
    t.balance_sum,
    t.total_supply,
    t.total_supply - t.balance_sum `offset`
FROM
    (SELECT
         addr_asset.asset_id,
         nep5.symbol,
         sum(addr_asset.balance) / 100000000 balance_sum,
         nep5.total_supply / 100000000       total_supply
     FROM addr_asset
         JOIN nep5
             on addr_asset.asset_id = nep5.asset_id
     GROUP BY nep5.asset_id) t
WHERE t.balance_sum != t.total_supply;

/* Verify total supply of asset */
SELECT
    t.asset_id,
    t.name,
    t.balance_sum,
    t.available,
    t.available - t.balance_sum `offset`
FROM
    (SELECT
         addr_asset.asset_id,
         asset.name,
         sum(addr_asset.balance) / 100000000 balance_sum,
         asset.available / 100000000         available
     FROM addr_asset
         JOIN asset
             on addr_asset.asset_id = asset.asset_id
     GROUP BY asset.asset_id) t
WHERE t.balance_sum != t.available;

/* Verify total count of different transaction types */
select 'cnt_tx_claim not equal' from counter where cnt_tx_claim != (select count(id) from tx where type='ClaimTransaction') union all
select 'cnt_tx_contract not equal' from counter where cnt_tx_contract != (select count(id) from tx where type='ContractTransaction') union all
select 'cnt_tx_invocation not equal' from counter where cnt_tx_invocation != (select count(id) from tx where type='InvocationTransaction') union all
select 'cnt_tx_issue not equal' from counter where cnt_tx_issue != (select count(id) from tx where type='IssueTransaction') union all
select 'cnt_tx_miner not equal' from counter where cnt_tx_miner != (select count(id) from tx where type='MinerTransaction') union all
select 'cnt_tx_reg not equal' from counter where cnt_tx_reg != (select count(id) from tx where type='RegisterTransaction') union all
select 'cnt_tx_publish not equal' from counter where cnt_tx_publish != (select count(id) from tx where type='PublishTransaction') union all
select 'cnt_tx_enrollment not equal' from counter where cnt_tx_enrollment != (select count(id) from tx where type='EnrollmentTransaction');

/*
    Check if address is missing but has addr_asset records,
    The result should always be empty
*/
SELECT ads.address, `a`.`address`, `a`.`created_at`, `a`.`last_transaction_time`, `ads`.`asset_id`, `ads`.`balance` FROM `addr_asset` as ads LEFT JOIN `address` a ON `a`.`address`=`ads`.`address` where a.address is null order by ads.id asc limit 10;


/* ---------------------------------------------------------------------- */

/* Get GAS available */
SELECT sum(tx_vout.value)
FROM tx_vout
    JOIN tx ON tx_vout.txid = tx.txid
WHERE tx.type = 'ClaimTransaction' AND
      tx_vout.asset_id = '0x602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7';

/* Get asset available */
SELECT
    asset.name,
    asset.asset_id,
    sum(tx_vout.value)
FROM tx_vout
    JOIN asset
    JOIN tx ON tx_vout.asset_id = asset.asset_id AND tx_vout.txid = tx.txid
WHERE tx.type = 'IssueTransaction'
      AND asset.asset_id NOT IN (
    '0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b',
    '0x602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7'
)
GROUP BY tx_vout.asset_id
ORDER BY asset.id DESC;