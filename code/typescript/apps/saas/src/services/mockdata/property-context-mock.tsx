import mlsData from './mls_listings_with_parcels.json';
import {
  convertMlsDataToPropertyDetails,
  convertMlsDataArrayToPropertyDetails,
} from './mls-data-converter';
import { PropertyDetails } from '@am/commons/schemas';

/**
 * Utility functions for working with MLS data
 */

/**
 * Filter an array of MLS records to keep only the latest record for each unique ID
 * based on the MLSRecordID (assuming higher MLSRecordID values are more recent)
 *
 * @param records The array of MLS records that may contain duplicates
 * @returns A new array with only the latest record for each unique ID
 */
export function filterDuplicateMlsRecords<T extends { ATTOM_ID: string; MLSRecordID: number }>(
  records: T[]
): T[] {
  const recordMap = new Map<string, T>();

  // Iterate through all records
  for (const record of records) {
    const existingRecord = recordMap.get(record.ATTOM_ID);

    // If we've never seen this ID before or if this record has a higher MLSRecordID,
    // update the map with this record
    if (!existingRecord || record.MLSRecordID > existingRecord.MLSRecordID) {
      recordMap.set(record.ATTOM_ID, record);
    }
  }

  // Convert the Map back to an array
  return Array.from(recordMap.values());
}

// Define the comparable property type based on PropertyDetails
export interface ComparableProperty extends PropertyDetails {
  distance?: number;
  closeDate?: string;
  amount?: string;
  mapboxData?: any;
  selected?: boolean;
}

// Take all properties from the MLS data
const mlsProperties = mlsData as any[];

// Filter duplicate records, keeping only the latest MLSRecordID for each ID
const uniqueMlsProperties = filterDuplicateMlsRecords(mlsProperties);

// Convert MLS data to PropertyDetails format for the primary property
export const mockPropertyData: PropertyDetails = convertMlsDataToPropertyDetails(
  uniqueMlsProperties[0]
);

// Convert additional properties to comparable properties (up to 1000)
export const mockComparableProperties: ComparableProperty[] = convertMlsDataArrayToPropertyDetails(
  uniqueMlsProperties.slice(1)
);
