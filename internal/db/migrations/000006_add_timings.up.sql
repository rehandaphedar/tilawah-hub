ALTER TABLE recitation_files
ADD COLUMN has_timings BOOLEAN NOT NULL DEFAULT 0;

ALTER TABLE recitation_files
ADD COLUMN lafzize_processing BOOLEAN NOT NULL DEFAULT 0;
