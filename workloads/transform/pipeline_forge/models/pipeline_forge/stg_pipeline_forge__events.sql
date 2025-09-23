{{ config(materialized='view') }}

with source_data as (

    select * from {{ source('pipeline_forge', 'events') }}

),

renamed as (

    select
        id as event_id,
        event as event_name
        
    from source_data

)

select * from renamed
