const axios = require("axios");
const { env } = require("../../config");

const handleGetStudentReport = async (req, res) => {
    const { id } = req.params;
    const goServiceUrl = env.GO_SERVICE_URL;

    try {
        const response = await axios.get(`${goServiceUrl}/api/v1/students/${id}/report`, {
            responseType: "stream",
            timeout: 10000,
        });

        res.setHeader("Content-Type", "application/pdf");
        res.setHeader(
            "Content-Disposition",
            response.headers["content-disposition"] || `attachment; filename=student_${id}_report.pdf`
        );

        response.data.pipe(res);
    } catch (error) {
        if (error.response?.status === 404) {
            return res.status(404).json({ message: "Student not found" });
        }
        if (error.code === "ECONNREFUSED" || error.code === "ETIMEDOUT") {
            return res.status(503).json({ message: "PDF service unavailable" });
        }
        return res.status(500).json({ message: "Failed to generate report" });
    }
};

module.exports = { handleGetStudentReport };
