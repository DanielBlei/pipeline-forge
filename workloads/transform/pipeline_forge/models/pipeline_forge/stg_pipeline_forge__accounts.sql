{{ config(materialized='view') }}

with source_data as (

    select * from {{ source('pipeline_forge', 'account') }}

),

renamed as (

    select
        id as account_id,
        email as account_email,

    from source_data

)

select * from renamed
