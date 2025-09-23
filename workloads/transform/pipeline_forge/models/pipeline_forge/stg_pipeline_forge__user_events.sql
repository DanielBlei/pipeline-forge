{{ config(materialized='view') }}

with source_data as (

    select * from {{ source('pipeline_forge', 'user_events') }}

),

renamed as (

    select
        account_id,
        event_id

    from source_data

)

select * from renamed
