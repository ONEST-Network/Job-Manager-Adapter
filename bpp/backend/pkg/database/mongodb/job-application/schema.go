package jobapplication

type JobApplication struct {
	ID          string `bson:"id"`
	JobID       string `bson:"job_id"`
	ApplicantID string `bson:"applicant_id"`
}
