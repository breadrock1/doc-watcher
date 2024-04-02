#!/bin/bash

NOTIFIER_INDEXER_DIR=$PWD/indexer
SEARCHER_SERVICE_ADDRESS=http://localhost:2892

TEST_BUCKET_NAME=common_bucket

TESTCASE_DIR_PATH=$PWD/testcases/integration
TESTCASE_FILE_1_HASH=fd48ca99fa980b47b100b3afb9d03b4b
TESTCASE_FILE_2_HASH=c2970f52b60198b49e3ef51101ed418b

function print_info_message() {
    echo -e "$(date) - INFO: $1"
}

function print_err_message() {
    echo -e "$(date) - WARN: $1"
}

function clear_all_data() {
    rm -rf $NOTIFIER_INDEXER_DIR/*

    curl -X DELETE --silent \
                   --url "$SEARCHER_SERVICE_ADDRESS/bucket/$TEST_BUCKET_NAME" \
                   --header 'Content-Type: application/json'
}

function create_integration_test_bucket() {
    RESPONSE=$(curl -X POST --silent \
                            --url "$SEARCHER_SERVICE_ADDRESS/bucket/new" \
                            --header 'Content-Type: application/json' \
                            --data '{"bucket_name": "common_bucket"}')

    return $(echo $RESPONSE | jq '.code')
}

function get_existing_buckets() {
    RESPONSE=$(curl -X GET --silent \
                           --url "$SEARCHER_SERVICE_ADDRESS/bucket/all" \
                           --header 'Content-Type: application/json')

    return $( echo $? )
}

function check_document_stored_successful() {
    TARGET_URL="$SEARCHER_SERVICE_ADDRESS/document/$TEST_BUCKET_NAME/$(md5 -q $1)"
    RESPONSE=$(curl -X GET --silent \
                           --url $TARGET_URL \
                           --header 'Content-Type: application/json')

    return $( echo $? )
}

function search_document_data() {
    JSON_QUERY=$(jq -n --arg query "$1" \
                       --arg buckets "$TEST_BUCKET_NAME" \
                       --arg document_type 'document' \
                       --arg scroll_timelife '1m' \
                       --arg document_extension '' \
                       --arg created_date_to '' \
                       --arg created_date_from '' \
                       --argjson result_size 5  \
                       --argjson result_offset 0  \
                       --argjson document_size_to 0  \
                       --argjson document_size_from 0  \
                       '$ARGS.named')

    RESPONSE=$(curl -X POST --silent \
                 --url "$SEARCHER_SERVICE_ADDRESS/search/" \
                 --header 'Content-Type: application/json' \
                 --data "$JSON_QUERY")

    return $( echo $? )
}

print_info_message "Creating test bucket into search service..."
RESPONSE=$(create_integration_test_bucket)
if [[ $RESPONSE == 200 ]]; then
    print_err_message "Failed while creating test bucket!"
    clear_all_data
    exit -1
fi

print_info_message "Checking that created bucket does exist..."
EXEC_CODE=$(get_existing_buckets)
GET_TEST_BUCKET_NAME=$(echo $RESPONSE | jq '.[].index' | grep $TEST_BUCKET_NAME)
if [[ EXEC_CODE != 0 && -z GET_TEST_BUCKET_NAME ]]; then
    print_err_message "There is no bucket with name $TEST_BUCKET_NAME"
    clear_all_data
    exit -1
fi

print_info_message "Coping testcase files '$TESTCASE_DIR_PATH' to indexer directory..."
EXEC_CODE=$(cp -r $TESTCASE_DIR_PATH $NOTIFIER_INDEXER_DIR/)
if [[ $( echo $? ) != 0 ]]; then
    print_err_message "Failed while coping testcase file to indexer!"
    clear_all_data
    exit -1
fi

print_info_message "Waiting storing documents to elastic..." && sleep 10

print_info_message "Check documents are stored successful into elastic..."
for TESTCASE_FILE in $(ls $TESTCASE_DIR_PATH); do
    EXEC_CODE=$(check_document_stored_successful $TESTCASE_DIR_PATH/$TESTCASE_FILE)
    EXEC_STATUS=$(echo $RESPONSE | jq '.document_md5')
    if [[ $( echo $? ) != 0 && -z $EXEC_STATUS ]]; then
        print_err_message "Failed while storing document to elastic!"
        clear_all_data
        exit -1
    fi
done

print_info_message "Searching data by loaded documents..."
EXEC_CODE=$(search_document_data 'System of indexing and storing documents')
FOUNDED_DOC_1=$(cat /tmp/result_output.json | jq ".founded.$TESTCASE_FILE_1_HASH.[0].document_md5")
if [[ -z $FOUNDED_DOC_1 && $FOUNDED_DOC_1 == $TESTCASE_FILE_1_HASH ]]; then
    print_err_message "Failed while searching document!"
    clear_all_data
    exit -1
fi

EXEC_CODE=$(search_document_data 'Quickly find documents based on content')
FOUNDED_DOC_2=$(cat /tmp/result_output.json | jq ".founded.$TESTCASE_FILE_2_HASH.[0].document_md5")
if [[ -z $FOUNDED_DOC_2 && $$FOUNDED_DOC_2 == $TESTCASE_FILE_2_HASH ]]; then
    print_err_message "Failed while searching document!"
    clear_all_data
    exit -1
fi

clear_all_data
