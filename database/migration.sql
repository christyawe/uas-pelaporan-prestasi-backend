ALTER TABLE achievement_references 
DROP CONSTRAINT IF EXISTS achievement_references_status_check;

ALTER TABLE achievement_references 
ADD CONSTRAINT achievement_references_status_check 
CHECK (status IN ('draft', 'submitted', 'verified', 'rejected', 'deleted'));